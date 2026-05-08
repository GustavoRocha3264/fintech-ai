import re
from dataclasses import dataclass
from pathlib import Path

from app.core.wiki_paths import CATEGORIES, category_of, is_markdown, resolve_within
from app.schemas.ingest import WikiEdit

FRONT_MATTER_RE = re.compile(r"^---\n(?P<body>.*?)\n---\n", re.DOTALL)
WIKILINK_RE = re.compile(r"\[\[([^\]]+)\]\]")
REQUIRED_KEYS = ("type", "status", "last_reviewed")
ALLOWED_STATUS = {"draft", "stable"}


@dataclass(frozen=True)
class ValidationError:
    path: str
    message: str


class IngestValidator:
    """Validates a batch of edits before any disk writes happen."""

    def __init__(self, *, wiki_root: Path) -> None:
        self._root = wiki_root

    def validate(self, edits: list[WikiEdit]) -> list[ValidationError]:
        errors: list[ValidationError] = []
        seen: set[str] = set()
        for edit in edits:
            if edit.path in seen:
                errors.append(ValidationError(edit.path, "duplicate path in batch"))
                continue
            seen.add(edit.path)
            errors.extend(self._validate_one(edit))
        return errors

    def _validate_one(self, edit: WikiEdit) -> list[ValidationError]:
        errs: list[ValidationError] = []
        path = edit.path

        if not is_markdown(path):
            errs.append(ValidationError(path, "must end in .md"))
        if category_of(path) is None:
            errs.append(ValidationError(path, f"first segment must be one of {CATEGORIES}"))
        try:
            resolved = resolve_within(self._root, path)
        except ValueError as e:
            errs.append(ValidationError(path, str(e)))
            return errs

        meta = _parse_front_matter(edit.content)
        if meta is None:
            errs.append(ValidationError(path, "missing YAML front-matter block"))
        else:
            for key in REQUIRED_KEYS:
                if key not in meta:
                    errs.append(ValidationError(path, f"front-matter missing '{key}'"))
            if meta.get("status") and meta["status"] not in ALLOWED_STATUS:
                errs.append(
                    ValidationError(path, f"status must be one of {sorted(ALLOWED_STATUS)}")
                )

        exists = resolved.is_file()
        if edit.op == "create" and exists:
            errs.append(ValidationError(path, "create requested but file already exists"))
        if edit.op == "update" and not exists:
            errs.append(ValidationError(path, "update requested but file is missing"))
        if edit.op == "update" and exists and not edit.force:
            existing = _parse_front_matter(resolved.read_text(encoding="utf-8")) or {}
            if existing.get("status") == "stable":
                errs.append(
                    ValidationError(path, "refusing to overwrite stable page without force=true")
                )

        return errs


def _parse_front_matter(text: str) -> dict[str, str] | None:
    m = FRONT_MATTER_RE.match(text)
    if not m:
        return None
    out: dict[str, str] = {}
    for raw in m.group("body").splitlines():
        line = raw.strip()
        if not line or line.startswith("#"):
            continue
        if ":" not in line:
            continue
        key, _, value = line.partition(":")
        out[key.strip()] = value.strip().strip('"').strip("'")
    return out


def extract_wikilinks(text: str) -> list[str]:
    return [m.group(1).strip() for m in WIKILINK_RE.finditer(text)]

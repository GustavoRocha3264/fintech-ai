from dataclasses import dataclass
from pathlib import Path

from app.core.wiki_paths import CATEGORIES, INDEX_FILE, SCHEMA_FILE, resolve_within


@dataclass(frozen=True)
class WikiPage:
    rel_path: str
    title: str
    body: str


class WikiReader:
    """Read-only view over the on-disk wiki. Index-first retrieval — the index
    is the entire content map for small wikis (Karpathy: stop-gap until ~300 pages).
    """

    def __init__(self, root: Path) -> None:
        self._root = root

    def root(self) -> Path:
        return self._root

    def exists(self) -> bool:
        return self._root.is_dir()

    def read_index(self) -> str:
        path = self._root / INDEX_FILE
        if not path.is_file():
            return ""
        return path.read_text(encoding="utf-8")

    def read_schema(self) -> str:
        path = self._root / SCHEMA_FILE
        if not path.is_file():
            return ""
        return path.read_text(encoding="utf-8")

    def read_page(self, rel_path: str) -> WikiPage | None:
        try:
            path = resolve_within(self._root, rel_path)
        except ValueError:
            return None
        if not path.is_file():
            return None
        body = path.read_text(encoding="utf-8")
        title = _extract_title(body) or path.stem
        return WikiPage(rel_path=rel_path, title=title, body=body)

    def list_pages(self) -> list[str]:
        out: list[str] = []
        for category in CATEGORIES:
            cat_dir = self._root / category
            if not cat_dir.is_dir():
                continue
            for path in sorted(cat_dir.glob("*.md")):
                out.append(f"{category}/{path.name}")
        return out


def _extract_title(body: str) -> str | None:
    for line in body.splitlines():
        line = line.strip()
        if line.startswith("# "):
            return line[2:].strip()
    return None

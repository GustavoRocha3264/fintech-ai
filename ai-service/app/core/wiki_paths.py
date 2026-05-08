from pathlib import Path

CATEGORIES: tuple[str, ...] = (
    "entities",
    "concepts",
    "sources",
    "syntheses",
    "playbooks",
)

INDEX_FILE = "index.md"
LOG_FILE = "log.md"
SCHEMA_FILE = "CLAUDE.md"


def resolve_within(root: Path, rel: str) -> Path:
    """Resolve `rel` against `root`, refusing escapes via `..` or absolute paths.

    Returns the resolved absolute path. Raises ValueError on traversal.
    """
    if rel.startswith("/") or "\\" in rel:
        raise ValueError(f"invalid wiki path: {rel}")
    candidate = (root / rel).resolve()
    root_resolved = root.resolve()
    try:
        candidate.relative_to(root_resolved)
    except ValueError as err:
        raise ValueError(f"path escapes wiki root: {rel}") from err
    return candidate


def is_markdown(rel: str) -> bool:
    return rel.endswith(".md")


def category_of(rel: str) -> str | None:
    head = rel.split("/", 1)[0]
    return head if head in CATEGORIES else None

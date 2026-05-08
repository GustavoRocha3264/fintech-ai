"""Standalone wiki linter. Run as `python scripts/wiki_check.py [wiki_dir]`.

Exits non-zero if any page is missing front-matter or has a broken wikilink.
Designed for CI; mirrors the logic served by GET /v1/wiki/lint.
"""
from __future__ import annotations

import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent.parent))

from app.services.ingest_pipeline import _parse_front_matter, extract_wikilinks
from app.services.wiki_reader import WikiReader


def main(argv: list[str]) -> int:
    root = Path(argv[1]) if len(argv) > 1 else Path(__file__).resolve().parent.parent / "wiki"
    reader = WikiReader(root)
    if not reader.exists():
        print(f"wiki not found: {root}", file=sys.stderr)
        return 2

    pages = reader.list_pages()
    page_set = {p.removesuffix(".md") for p in pages}
    issues: list[str] = []

    for rel in pages:
        page = reader.read_page(rel)
        if page is None:
            issues.append(f"{rel}: unreadable")
            continue
        meta = _parse_front_matter(page.body)
        if meta is None:
            issues.append(f"{rel}: missing front-matter")
        else:
            for key in ("type", "status", "last_reviewed"):
                if key not in meta:
                    issues.append(f"{rel}: front-matter missing '{key}'")
        for link in extract_wikilinks(page.body):
            if link not in page_set:
                issues.append(f"{rel}: broken wikilink [[{link}]]")

    if issues:
        print(f"wiki lint: {len(issues)} issue(s) in {len(pages)} pages", file=sys.stderr)
        for line in issues:
            print(f"  - {line}", file=sys.stderr)
        return 1

    print(f"wiki lint: OK ({len(pages)} pages)")
    return 0


if __name__ == "__main__":
    raise SystemExit(main(sys.argv))

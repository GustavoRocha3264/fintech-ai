from __future__ import annotations

from pathlib import Path

import pytest
from fastapi.testclient import TestClient

from app.api import dependencies as deps
from app.core.config import Settings, get_settings
from app.main import app
from app.services.ingest_pipeline import IngestValidator, _parse_front_matter, extract_wikilinks
from app.services.wiki_reader import WikiReader
from app.services.wiki_writer import WikiWriter


client = TestClient(app)


@pytest.fixture
def tmp_wiki(tmp_path: Path) -> Path:
    for cat in ("entities", "concepts", "sources", "syntheses", "playbooks"):
        (tmp_path / cat).mkdir(parents=True)
    (tmp_path / "CLAUDE.md").write_text("# schema\n", encoding="utf-8")
    (tmp_path / "index.md").write_text("# index\n", encoding="utf-8")
    (tmp_path / "log.md").write_text("# log\n", encoding="utf-8")
    return tmp_path


@pytest.fixture
def override_wiki(tmp_wiki: Path):
    settings = Settings(
        anthropic_api_key="test",
        wiki_dir=tmp_wiki,
        wiki_ingest_token="secret",
    )
    app.dependency_overrides[get_settings] = lambda: settings
    app.dependency_overrides[deps.get_wiki_reader] = lambda: WikiReader(tmp_wiki)
    app.dependency_overrides[deps.get_wiki_writer] = lambda: WikiWriter(wiki_root=tmp_wiki)
    app.dependency_overrides[deps.get_ingest_validator] = lambda: IngestValidator(wiki_root=tmp_wiki)
    yield tmp_wiki
    for key in (
        get_settings,
        deps.get_wiki_reader,
        deps.get_wiki_writer,
        deps.get_ingest_validator,
    ):
        app.dependency_overrides.pop(key, None)


def _page(slug: str, status: str = "draft") -> str:
    return (
        "---\n"
        "type: concept\n"
        f"status: {status}\n"
        "last_reviewed: 2026-05-08\n"
        "---\n\n"
        f"# {slug}\n"
    )


def test_wiki_index_lists_real_pages() -> None:
    r = client.get("/v1/wiki/index")
    assert r.status_code == 200
    pages = r.json()["pages"]
    assert "entities/USD.md" in pages
    assert "concepts/PFIC.md" in pages


def test_wiki_lint_passes_on_seed_wiki() -> None:
    r = client.get("/v1/wiki/lint")
    assert r.status_code == 200
    assert r.json()["issues"] == []


def test_ingest_requires_bearer_token(override_wiki: Path) -> None:
    payload = {
        "source": {"name": "x", "fetched_at": "2026-05-08T00:00:00Z"},
        "edits": [{"op": "create", "path": "concepts/foo.md", "content": _page("foo")}],
        "log_entry": "test",
    }
    r = client.post("/v1/wiki/ingest", json=payload)
    assert r.status_code == 401


def test_ingest_creates_page(override_wiki: Path) -> None:
    payload = {
        "source": {"name": "test", "fetched_at": "2026-05-08T00:00:00Z"},
        "edits": [{"op": "create", "path": "concepts/foo.md", "content": _page("foo")}],
        "log_entry": "first foo",
    }
    r = client.post(
        "/v1/wiki/ingest",
        json=payload,
        headers={"Authorization": "Bearer secret"},
    )
    assert r.status_code == 200, r.text
    assert (override_wiki / "concepts" / "foo.md").is_file()
    log_text = (override_wiki / "log.md").read_text(encoding="utf-8")
    assert "first foo" in log_text


def test_ingest_rejects_path_traversal(override_wiki: Path) -> None:
    payload = {
        "source": {"name": "evil", "fetched_at": "2026-05-08T00:00:00Z"},
        "edits": [{"op": "create", "path": "../etc.md", "content": _page("evil")}],
        "log_entry": "x",
    }
    r = client.post(
        "/v1/wiki/ingest",
        json=payload,
        headers={"Authorization": "Bearer secret"},
    )
    assert r.status_code == 422


def test_ingest_rejects_clobber_of_stable_without_force(override_wiki: Path) -> None:
    target = override_wiki / "concepts" / "stable-thing.md"
    target.write_text(_page("stable-thing", status="stable"), encoding="utf-8")
    payload = {
        "source": {"name": "weekly", "fetched_at": "2026-05-08T00:00:00Z"},
        "edits": [
            {
                "op": "update",
                "path": "concepts/stable-thing.md",
                "content": _page("stable-thing", status="stable"),
            }
        ],
        "log_entry": "trying to overwrite",
    }
    r = client.post(
        "/v1/wiki/ingest",
        json=payload,
        headers={"Authorization": "Bearer secret"},
    )
    assert r.status_code == 422
    assert "stable" in r.text


def test_ingest_atomic_rollback_on_failure(tmp_wiki: Path) -> None:
    (tmp_wiki / "concepts" / "a.md").write_text(_page("a"), encoding="utf-8")
    writer = WikiWriter(wiki_root=tmp_wiki)
    from app.schemas.ingest import IngestRequest, IngestSource, WikiEdit

    req = IngestRequest(
        source=IngestSource(name="t", fetched_at="2026-05-08T00:00:00Z"),
        edits=[
            WikiEdit(op="update", path="concepts/a.md", content=_page("a-v2")),
            WikiEdit(op="create", path="concepts/b.md", content=_page("b")),
        ],
        log_entry="batch",
    )

    import os as _os

    original_replace = _os.replace
    calls = {"n": 0}

    def flaky_replace(src, dst):
        calls["n"] += 1
        if calls["n"] == 3:
            raise OSError("simulated failure")
        return original_replace(src, dst)

    _os.replace = flaky_replace
    try:
        with pytest.raises(OSError):
            writer.apply(req)
    finally:
        _os.replace = original_replace

    assert (tmp_wiki / "concepts" / "a.md").read_text(encoding="utf-8") == _page("a")
    assert not (tmp_wiki / "concepts" / "b.md").exists()


def test_front_matter_parser() -> None:
    meta = _parse_front_matter(_page("x"))
    assert meta is not None
    assert meta["type"] == "concept"
    assert meta["status"] == "draft"


def test_extract_wikilinks() -> None:
    text = "see [[entities/USD]] and [[concepts/PFIC]]."
    assert extract_wikilinks(text) == ["entities/USD", "concepts/PFIC"]

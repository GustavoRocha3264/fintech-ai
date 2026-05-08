from fastapi.testclient import TestClient

from app.main import app

client = TestClient(app)


def test_chat_grounds_in_wiki(fake_llm) -> None:
    fake_llm.reply = "exposição cambial é..."
    r = client.post(
        "/v1/chat",
        json={"messages": [{"role": "user", "content": "fala sobre exposicao-cambial"}]},
    )
    assert r.status_code == 200
    body = r.json()
    assert body["answer"] == "exposição cambial é..."
    assert "concepts/exposicao-cambial.md" in body["citations"]
    assert "Wiki schema" in (fake_llm.captured_system or "")
    assert "Wiki index" in (fake_llm.captured_system or "")


def test_chat_rejects_empty_messages() -> None:
    r = client.post("/v1/chat", json={"messages": []})
    assert r.status_code == 422


def test_chat_rejects_invalid_role() -> None:
    r = client.post(
        "/v1/chat",
        json={"messages": [{"role": "system", "content": "hi"}]},
    )
    assert r.status_code == 422


def test_reports_endpoint_is_deprecated() -> None:
    r = client.post("/v1/reports", json={"portfolio_id": "p1", "holdings": []})
    assert r.status_code == 200
    assert r.headers.get("Deprecation") == "true"
    assert "successor-version" in r.headers.get("Link", "")

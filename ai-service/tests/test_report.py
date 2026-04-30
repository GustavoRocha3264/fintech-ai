from fastapi.testclient import TestClient

from app.main import app

client = TestClient(app)


def test_healthz() -> None:
    r = client.get("/healthz")
    assert r.status_code == 200
    assert r.json() == {"status": "ok"}


def test_generate_report_stub() -> None:
    r = client.post("/v1/reports", json={"portfolio_id": "p1", "holdings": []})
    assert r.status_code == 200
    body = r.json()
    assert body["portfolio_id"] == "p1"
    assert "risk" in body


def test_reports_endpoint_advertises_deprecation() -> None:
    r = client.post("/v1/reports", json={"portfolio_id": "p1", "holdings": []})
    assert r.headers.get("Deprecation") == "true"
    assert "/v1/chat" in r.headers.get("Link", "")

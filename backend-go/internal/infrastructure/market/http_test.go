package market_test

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/fintech/cbpi/backend-go/internal/infrastructure/market"
)

func TestHTTPProvider_FetchesAndCaches(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"results":[{"symbol":"PETR4","regularMarketPrice":38.50,"currency":"BRL"}]}`))
	}))
	defer srv.Close()

	p := market.NewHTTPProvider(srv.URL, time.Minute, "")

	price, currency, err := p.GetPrice("PETR4")
	if err != nil {
		t.Fatalf("GetPrice: %v", err)
	}
	if price != 38.50 {
		t.Fatalf("price = %v, want 38.50", price)
	}
	if currency != "BRL" {
		t.Fatalf("currency = %v, want BRL", currency)
	}

	// Second call within TTL must be served from cache — no extra HTTP hit.
	if _, _, err := p.GetPrice("PETR4"); err != nil {
		t.Fatalf("second GetPrice: %v", err)
	}
	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Fatalf("expected 1 HTTP call (cached after first), got %d", got)
	}
}

func TestHTTPProvider_USDSymbol(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"results":[{"symbol":"AAPL","regularMarketPrice":195.40,"currency":"USD"}]}`))
	}))
	defer srv.Close()

	p := market.NewHTTPProvider(srv.URL, time.Minute, "")
	price, currency, err := p.GetPrice("AAPL")
	if err != nil {
		t.Fatalf("GetPrice: %v", err)
	}
	if price != 195.40 {
		t.Fatalf("price = %v, want 195.40", price)
	}
	if currency != "USD" {
		t.Fatalf("currency = %v, want USD", currency)
	}
}

func TestHTTPProvider_RespectsTTL(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		atomic.AddInt32(&calls, 1)
		_, _ = w.Write([]byte(`{"results":[{"symbol":"VALE3","regularMarketPrice":65.20,"currency":"BRL"}]}`))
	}))
	defer srv.Close()

	p := market.NewHTTPProvider(srv.URL, 1*time.Millisecond, "")
	if _, _, err := p.GetPrice("VALE3"); err != nil {
		t.Fatal(err)
	}
	time.Sleep(5 * time.Millisecond)
	if _, _, err := p.GetPrice("VALE3"); err != nil {
		t.Fatal(err)
	}
	if got := atomic.LoadInt32(&calls); got != 2 {
		t.Fatalf("expected 2 calls after TTL expiry, got %d", got)
	}
}

func TestHTTPProvider_ErrorOnBadStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer srv.Close()

	p := market.NewHTTPProvider(srv.URL, time.Minute, "")
	if _, _, err := p.GetPrice("PETR4"); err == nil {
		t.Fatal("expected error on 500")
	}
}

func TestHTTPProvider_ErrorOnAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"error":true,"message":"Invalid symbol: BOGUS"}`))
	}))
	defer srv.Close()

	p := market.NewHTTPProvider(srv.URL, time.Minute, "")
	if _, _, err := p.GetPrice("BOGUS"); err == nil {
		t.Fatal("expected error for API error response")
	}
}

func TestHTTPProvider_ErrorOnEmptyResults(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"results":[]}`))
	}))
	defer srv.Close()

	p := market.NewHTTPProvider(srv.URL, time.Minute, "")
	if _, _, err := p.GetPrice("PETR4"); err == nil {
		t.Fatal("expected error when results array is empty")
	}
}

func TestHTTPProvider_TokenPassedAsQueryParam(t *testing.T) {
	var gotToken string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotToken = r.URL.Query().Get("token")
		_, _ = w.Write([]byte(`{"results":[{"symbol":"PETR4","regularMarketPrice":38.50,"currency":"BRL"}]}`))
	}))
	defer srv.Close()

	p := market.NewHTTPProvider(srv.URL, time.Minute, "my-secret-token")
	if _, _, err := p.GetPrice("PETR4"); err != nil {
		t.Fatalf("GetPrice: %v", err)
	}
	if gotToken != "my-secret-token" {
		t.Fatalf("token = %q, want %q", gotToken, "my-secret-token")
	}
}

func TestHTTPProvider_NoTokenQueryParamWhenEmpty(t *testing.T) {
	var gotURL string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotURL = r.URL.RawQuery
		_, _ = w.Write([]byte(`{"results":[{"symbol":"PETR4","regularMarketPrice":38.50,"currency":"BRL"}]}`))
	}))
	defer srv.Close()

	p := market.NewHTTPProvider(srv.URL, time.Minute, "")
	if _, _, err := p.GetPrice("PETR4"); err != nil {
		t.Fatalf("GetPrice: %v", err)
	}
	if gotURL != "" {
		t.Fatalf("expected no query string when token is empty, got %q", gotURL)
	}
}

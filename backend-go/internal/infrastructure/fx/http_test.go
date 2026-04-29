package fx_test

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/fintech/cbpi/backend-go/internal/infrastructure/fx"
)

func TestHTTPProvider_FetchesAndCaches(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"result":"success","base_code":"USD","rates":{"USD":1,"BRL":5.07,"EUR":0.92}}`))
	}))
	defer srv.Close()

	p := fx.NewHTTPProvider(srv.URL, time.Minute)

	r, err := p.GetRate("USD", "BRL")
	if err != nil {
		t.Fatalf("get USD->BRL: %v", err)
	}
	if r != 5.07 {
		t.Fatalf("USD->BRL = %v", r)
	}

	r, err = p.GetRate("USD", "EUR")
	if err != nil {
		t.Fatalf("get USD->EUR: %v", err)
	}
	if r != 0.92 {
		t.Fatalf("USD->EUR = %v", r)
	}

	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Fatalf("expected 1 HTTP call (cached after first), got %d", got)
	}
}

func TestHTTPProvider_SameCurrencyReturnsOne(t *testing.T) {
	p := fx.NewHTTPProvider("http://unused", time.Minute)
	r, err := p.GetRate("USD", "USD")
	if err != nil || r != 1.0 {
		t.Fatalf("expected 1.0/nil, got %v/%v", r, err)
	}
}

func TestHTTPProvider_RespectsTTL(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		atomic.AddInt32(&calls, 1)
		_, _ = w.Write([]byte(`{"result":"success","base_code":"USD","rates":{"BRL":5.0}}`))
	}))
	defer srv.Close()

	p := fx.NewHTTPProvider(srv.URL, 1*time.Millisecond)
	if _, err := p.GetRate("USD", "BRL"); err != nil {
		t.Fatal(err)
	}
	time.Sleep(5 * time.Millisecond)
	if _, err := p.GetRate("USD", "BRL"); err != nil {
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

	p := fx.NewHTTPProvider(srv.URL, time.Minute)
	if _, err := p.GetRate("USD", "BRL"); err == nil {
		t.Fatal("expected error on 500")
	}
}

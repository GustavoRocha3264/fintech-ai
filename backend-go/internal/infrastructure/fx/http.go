package fx

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// HTTPProvider fetches live FX rates from a public API and caches them.
// Default endpoint: https://open.er-api.com/v6/latest/{base} — free, no API
// key, supports BRL/USD. The response shape is:
//
//	{ "result":"success", "base_code":"USD", "rates":{"BRL":5.07,"USD":1,...} }
//
// One HTTP call per base currency populates the cache for every quote
// currency in the response, so subsequent GetRate calls within the TTL are
// in-memory.
type HTTPProvider struct {
	baseURL string
	client  *http.Client
	ttl     time.Duration

	mu    sync.RWMutex
	cache map[string]cachedRates // key: base currency
}

type cachedRates struct {
	rates     map[string]float64
	expiresAt time.Time
}

type apiResponse struct {
	Result  string             `json:"result"`
	Base    string             `json:"base_code"`
	Rates   map[string]float64 `json:"rates"`
	ErrType string             `json:"error-type"`
}

func NewHTTPProvider(baseURL string, ttl time.Duration) *HTTPProvider {
	return &HTTPProvider{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 5 * time.Second},
		ttl:     ttl,
		cache:   map[string]cachedRates{},
	}
}

func (p *HTTPProvider) GetRate(from, to string) (float64, error) {
	if from == to {
		return 1.0, nil
	}
	if r, ok := p.lookup(from, to); ok {
		return r, nil
	}
	if err := p.refresh(from); err != nil {
		return 0, err
	}
	if r, ok := p.lookup(from, to); ok {
		return r, nil
	}
	return 0, fmt.Errorf("rate not available: %s/%s", from, to)
}

func (p *HTTPProvider) lookup(from, to string) (float64, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	entry, ok := p.cache[from]
	if !ok || time.Now().After(entry.expiresAt) {
		return 0, false
	}
	r, ok := entry.rates[to]
	return r, ok
}

func (p *HTTPProvider) refresh(from string) error {
	url := fmt.Sprintf("%s/%s", p.baseURL, from)
	resp, err := p.client.Get(url)
	if err != nil {
		return fmt.Errorf("fx fetch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("fx fetch: status %d", resp.StatusCode)
	}

	var body apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return fmt.Errorf("fx decode: %w", err)
	}
	if body.Result != "" && body.Result != "success" {
		return fmt.Errorf("fx api error: %s", body.ErrType)
	}
	if len(body.Rates) == 0 {
		return fmt.Errorf("fx api returned no rates")
	}

	p.mu.Lock()
	p.cache[from] = cachedRates{rates: body.Rates, expiresAt: time.Now().Add(p.ttl)}
	p.mu.Unlock()
	return nil
}

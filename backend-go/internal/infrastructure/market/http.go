package market

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/fintech/cbpi/backend-go/internal/domain/portfolio"
)

// HTTPProvider fetches live market quotes from Brapi (https://brapi.dev) and
// caches them per-symbol for the configured TTL.
//
// Default endpoint: https://brapi.dev/api/quote/{symbol}
// An optional bearer token can be set via BRAPI_TOKEN; without it the free
// unauthenticated tier is used (rate-limited but sufficient for dev/small use).
//
// Response shape (subset):
//
//	{ "results": [{ "symbol":"PETR4", "regularMarketPrice":38.50, "currency":"BRL" }] }
type HTTPProvider struct {
	baseURL string
	token   string
	client  *http.Client
	ttl     time.Duration

	mu    sync.RWMutex
	cache map[string]cachedQuote
}

type cachedQuote struct {
	price     float64
	currency  string
	expiresAt time.Time
}

type brapiResponse struct {
	Results []struct {
		Symbol             string  `json:"symbol"`
		RegularMarketPrice float64 `json:"regularMarketPrice"`
		Currency           string  `json:"currency"`
	} `json:"results"`
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func NewHTTPProvider(baseURL string, ttl time.Duration, token string) *HTTPProvider {
	return &HTTPProvider{
		baseURL: baseURL,
		token:   token,
		client:  &http.Client{Timeout: 5 * time.Second},
		ttl:     ttl,
		cache:   map[string]cachedQuote{},
	}
}

func (p *HTTPProvider) GetPrice(symbol string) (float64, string, error) {
	if q, ok := p.lookup(symbol); ok {
		return q.price, q.currency, nil
	}
	if err := p.refresh(symbol); err != nil {
		return 0, "", err
	}
	if q, ok := p.lookup(symbol); ok {
		return q.price, q.currency, nil
	}
	return 0, "", fmt.Errorf("market: price unavailable for %s", symbol)
}

func (p *HTTPProvider) lookup(symbol string) (cachedQuote, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	q, ok := p.cache[symbol]
	if !ok || time.Now().After(q.expiresAt) {
		return cachedQuote{}, false
	}
	return q, true
}

func (p *HTTPProvider) refresh(symbol string) error {
	url := fmt.Sprintf("%s/%s", p.baseURL, symbol)
	if p.token != "" {
		url += "?token=" + p.token
	}

	resp, err := p.client.Get(url) //nolint:noctx
	if err != nil {
		return fmt.Errorf("market fetch %s: %w", symbol, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("market fetch %s: status %d", symbol, resp.StatusCode)
	}

	var body brapiResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return fmt.Errorf("market decode %s: %w", symbol, err)
	}
	if body.Error {
		return fmt.Errorf("market api error for %s: %s", symbol, body.Message)
	}
	if len(body.Results) == 0 {
		return fmt.Errorf("market api returned no results for %s", symbol)
	}

	result := body.Results[0]

	p.mu.Lock()
	p.cache[symbol] = cachedQuote{
		price:     result.RegularMarketPrice,
		currency:  normalizeCurrency(result.Currency),
		expiresAt: time.Now().Add(p.ttl),
	}
	p.mu.Unlock()
	return nil
}

// normalizeCurrency maps whatever the API returns to our domain constants.
// Brapi already uses "BRL" / "USD", but we guard against unexpected values.
func normalizeCurrency(raw string) string {
	switch raw {
	case portfolio.CurrencyBRL:
		return portfolio.CurrencyBRL
	case portfolio.CurrencyUSD:
		return portfolio.CurrencyUSD
	default:
		return portfolio.CurrencyUSD
	}
}

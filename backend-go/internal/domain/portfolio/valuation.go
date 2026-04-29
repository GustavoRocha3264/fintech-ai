package portfolio

type Valuation struct {
	TotalBRL     Money
	TotalUSD     Money
	PercentInBRL float64
	PercentInUSD float64
}

// Valuate converts every position into both currencies using current market
// prices and an FX rate, then aggregates totals and the BRL/USD split.
//
// prices  : symbol -> Money (market price in its native currency)
// fxRate  : 1 USD = fxRate BRL
func Valuate(positions []Position, prices map[string]Money, fxRate float64) Valuation {
	var totalBRL, totalUSD float64

	for _, p := range positions {
		price, ok := prices[p.Symbol]
		if !ok {
			price = Money{Amount: p.Price, Currency: p.Currency}
		}
		nativeValue := p.Quantity * price.Amount
		switch price.Currency {
		case CurrencyBRL:
			totalBRL += nativeValue
			if fxRate > 0 {
				totalUSD += nativeValue / fxRate
			}
		case CurrencyUSD:
			totalUSD += nativeValue
			totalBRL += nativeValue * fxRate
		}
	}

	v := Valuation{
		TotalBRL: NewMoney(totalBRL, CurrencyBRL),
		TotalUSD: NewMoney(totalUSD, CurrencyUSD),
	}
	if totalBRL > 0 {
		brlNative, usdNativeInBRL := nativeBreakdownInBRL(positions, prices, fxRate)
		denom := brlNative + usdNativeInBRL
		if denom > 0 {
			v.PercentInBRL = (brlNative / denom) * 100
			v.PercentInUSD = (usdNativeInBRL / denom) * 100
		}
	}
	return v
}

func nativeBreakdownInBRL(positions []Position, prices map[string]Money, fxRate float64) (float64, float64) {
	var brl, usdInBRL float64
	for _, p := range positions {
		price, ok := prices[p.Symbol]
		if !ok {
			price = Money{Amount: p.Price, Currency: p.Currency}
		}
		nv := p.Quantity * price.Amount
		switch price.Currency {
		case CurrencyBRL:
			brl += nv
		case CurrencyUSD:
			usdInBRL += nv * fxRate
		}
	}
	return brl, usdInBRL
}

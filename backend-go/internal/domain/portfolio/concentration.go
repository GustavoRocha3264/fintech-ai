package portfolio

// TopAssetConcentration returns the largest single position's share of the
// portfolio's total value, expressed as a percentage. Cross-currency positions
// are normalized to BRL using the supplied USD→BRL rate so the comparison is
// like-for-like.
func TopAssetConcentration(positions []Position, prices map[string]Money, usdToBRL float64) float64 {
	var total, top float64
	for _, pos := range positions {
		price, ok := prices[pos.Symbol]
		if !ok {
			price = NewMoney(pos.Price, pos.Currency)
		}
		v := pos.Quantity * price.Amount
		if price.Currency == CurrencyUSD {
			v *= usdToBRL
		}
		total += v
		if v > top {
			top = v
		}
	}
	if total == 0 {
		return 0
	}
	return (top / total) * 100
}

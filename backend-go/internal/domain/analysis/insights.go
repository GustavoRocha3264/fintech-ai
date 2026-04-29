package analysis

const (
	highCurrencyExposureThreshold = 70.0
	highConcentrationThreshold    = 50.0
	lowDiversificationThreshold   = 3
)

// GenerateInsights applies deterministic rules to a portfolio's metrics and
// returns user-facing observations. Adding a new rule means appending another
// pure function to this list.
func GenerateInsights(in Input) []string {
	rules := []func(Input) (string, bool){
		highBRLExposure,
		highUSDExposure,
		highConcentration,
		lowDiversification,
	}
	insights := make([]string, 0, len(rules))
	for _, r := range rules {
		if msg, ok := r(in); ok {
			insights = append(insights, msg)
		}
	}
	return insights
}

func highBRLExposure(in Input) (string, bool) {
	if in.BRLExposurePercent > highCurrencyExposureThreshold {
		return "High exposure to BRL. Consider increasing USD diversification.", true
	}
	return "", false
}

func highUSDExposure(in Input) (string, bool) {
	if in.USDExposurePercent > highCurrencyExposureThreshold {
		return "High exposure to USD. Consider balancing with local assets.", true
	}
	return "", false
}

func highConcentration(in Input) (string, bool) {
	if in.TopAssetConcentrationPercent > highConcentrationThreshold {
		return "Portfolio is highly concentrated in a single asset.", true
	}
	return "", false
}

func lowDiversification(in Input) (string, bool) {
	if in.PositionCount > 0 && in.PositionCount < lowDiversificationThreshold {
		return "Low diversification. Consider adding more assets.", true
	}
	return "", false
}

package analysis_test

import (
	"strings"
	"testing"

	"github.com/fintech/cbpi/backend-go/internal/domain/analysis"
)

func TestInsights_HighBRLExposure(t *testing.T) {
	got := analysis.GenerateInsights(analysis.Input{
		BRLExposurePercent: 80, USDExposurePercent: 20, PositionCount: 5,
	})
	if !contains(got, "High exposure to BRL") {
		t.Fatalf("missing BRL insight: %v", got)
	}
}

func TestInsights_HighUSDExposure(t *testing.T) {
	got := analysis.GenerateInsights(analysis.Input{
		BRLExposurePercent: 20, USDExposurePercent: 80, PositionCount: 5,
	})
	if !contains(got, "High exposure to USD") {
		t.Fatalf("missing USD insight: %v", got)
	}
}

func TestInsights_HighConcentration(t *testing.T) {
	got := analysis.GenerateInsights(analysis.Input{
		BRLExposurePercent: 50, USDExposurePercent: 50,
		TopAssetConcentrationPercent: 60, PositionCount: 5,
	})
	if !contains(got, "highly concentrated") {
		t.Fatalf("missing concentration insight: %v", got)
	}
}

func TestInsights_LowDiversification(t *testing.T) {
	got := analysis.GenerateInsights(analysis.Input{
		BRLExposurePercent: 50, USDExposurePercent: 50, PositionCount: 2,
	})
	if !contains(got, "Low diversification") {
		t.Fatalf("missing diversification insight: %v", got)
	}
}

func TestInsights_BalancedHasNoFlags(t *testing.T) {
	got := analysis.GenerateInsights(analysis.Input{
		BRLExposurePercent: 50, USDExposurePercent: 50,
		TopAssetConcentrationPercent: 20, PositionCount: 8,
	})
	if len(got) != 0 {
		t.Fatalf("expected no insights, got %v", got)
	}
}

func contains(haystack []string, needle string) bool {
	for _, h := range haystack {
		if strings.Contains(h, needle) {
			return true
		}
	}
	return false
}

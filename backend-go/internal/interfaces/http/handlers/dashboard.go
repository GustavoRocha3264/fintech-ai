package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/fintech/cbpi/backend-go/internal/application/dashboard"
	"github.com/fintech/cbpi/backend-go/internal/domain/portfolio"
	"github.com/fintech/cbpi/backend-go/internal/interfaces/http/dto"
)

type DashboardHandler struct {
	uc *dashboard.GetDashboard
}

func NewDashboardHandler(uc *dashboard.GetDashboard) *DashboardHandler {
	return &DashboardHandler{uc: uc}
}

type fxBlock struct {
	USDToBRL float64 `json:"usdToBRL"`
	BRLToUSD float64 `json:"brlToUSD"`
}

type snapshotItem struct {
	Timestamp     time.Time `json:"timestamp"`
	TotalValueBRL float64   `json:"totalValueBRL"`
	TotalValueUSD float64   `json:"totalValueUSD"`
}

type dashboardResponse struct {
	Portfolio    dto.PortfolioResponse `json:"portfolio"`
	Valuation    dto.ValuationResponse `json:"valuation"`
	LatestReport any                   `json:"latestReport"` // analysisResponse or null
	Snapshots    []snapshotItem        `json:"snapshots"`
	FX           fxBlock               `json:"fx"`
}

func (h *DashboardHandler) Get(c *gin.Context) {
	d, err := h.uc.Execute(c.Param("id"))
	if err != nil {
		if errors.Is(err, portfolio.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	snaps := make([]snapshotItem, 0, len(d.Snapshots))
	for _, s := range d.Snapshots {
		snaps = append(snaps, snapshotItem{
			Timestamp:     s.Timestamp,
			TotalValueBRL: s.TotalValueBRL,
			TotalValueUSD: s.TotalValueUSD,
		})
	}

	var report any
	if d.LatestReport != nil {
		report = analysisResponse{
			ID:                           d.LatestReport.ID,
			PortfolioID:                  d.LatestReport.PortfolioID,
			CreatedAt:                    d.LatestReport.CreatedAt,
			TotalValueBRL:                d.LatestReport.TotalValueBRL,
			TotalValueUSD:                d.LatestReport.TotalValueUSD,
			BRLExposurePercent:           d.LatestReport.BRLExposurePercent,
			USDExposurePercent:           d.LatestReport.USDExposurePercent,
			TopAssetConcentrationPercent: d.LatestReport.TopAssetConcentrationPercent,
			Insights:                     d.LatestReport.Insights,
		}
	}

	c.JSON(http.StatusOK, dashboardResponse{
		Portfolio: toPortfolioResponse(d.Portfolio),
		Valuation: dto.ValuationResponse{
			TotalBRL:     dto.MoneyResponse{Amount: d.Valuation.TotalBRL.Amount, Currency: d.Valuation.TotalBRL.Currency},
			TotalUSD:     dto.MoneyResponse{Amount: d.Valuation.TotalUSD.Amount, Currency: d.Valuation.TotalUSD.Currency},
			PercentInBRL: d.Valuation.PercentInBRL,
			PercentInUSD: d.Valuation.PercentInUSD,
		},
		LatestReport: report,
		Snapshots:    snaps,
		FX:           fxBlock{USDToBRL: d.USDToBRL, BRLToUSD: d.BRLToUSD},
	})
}

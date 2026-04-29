package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	apanalysis "github.com/fintech/cbpi/backend-go/internal/application/analysis"
	"github.com/fintech/cbpi/backend-go/internal/domain/analysis"
	"github.com/fintech/cbpi/backend-go/internal/domain/portfolio"
)

type AnalysisHandler struct {
	run    *apanalysis.RunAnalysis
	latest *apanalysis.GetLatestAnalysis
}

func NewAnalysisHandler(run *apanalysis.RunAnalysis, latest *apanalysis.GetLatestAnalysis) *AnalysisHandler {
	return &AnalysisHandler{run: run, latest: latest}
}

type analysisResponse struct {
	ID                           string    `json:"id"`
	PortfolioID                  string    `json:"portfolioId"`
	CreatedAt                    time.Time `json:"createdAt"`
	TotalValueBRL                float64   `json:"totalValueBRL"`
	TotalValueUSD                float64   `json:"totalValueUSD"`
	BRLExposurePercent           float64   `json:"brlExposurePercent"`
	USDExposurePercent           float64   `json:"usdExposurePercent"`
	TopAssetConcentrationPercent float64   `json:"topAssetConcentrationPercent"`
	Insights                     []string  `json:"insights"`
}

func (h *AnalysisHandler) Run(c *gin.Context) {
	report, err := h.run.Execute(c.Param("id"))
	if err != nil {
		if errors.Is(err, portfolio.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, toAnalysisResponse(report))
}

func (h *AnalysisHandler) Latest(c *gin.Context) {
	report, err := h.latest.Execute(c.Param("id"))
	if err != nil {
		if errors.Is(err, analysis.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, toAnalysisResponse(report))
}

func toAnalysisResponse(r *analysis.AnalysisReport) analysisResponse {
	return analysisResponse{
		ID:                           r.ID,
		PortfolioID:                  r.PortfolioID,
		CreatedAt:                    r.CreatedAt,
		TotalValueBRL:                r.TotalValueBRL,
		TotalValueUSD:                r.TotalValueUSD,
		BRLExposurePercent:           r.BRLExposurePercent,
		USDExposurePercent:           r.USDExposurePercent,
		TopAssetConcentrationPercent: r.TopAssetConcentrationPercent,
		Insights:                     r.Insights,
	}
}

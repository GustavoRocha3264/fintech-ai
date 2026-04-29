package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	apanalysis "github.com/fintech/cbpi/backend-go/internal/application/analysis"
	"github.com/fintech/cbpi/backend-go/internal/domain/analysis"
	"github.com/fintech/cbpi/backend-go/internal/domain/portfolio"
	"github.com/fintech/cbpi/backend-go/internal/interfaces/http/dto"
)

type AnalysisHandler struct {
	run    *apanalysis.RunAnalysis
	latest *apanalysis.GetLatestAnalysis
}

func NewAnalysisHandler(run *apanalysis.RunAnalysis, latest *apanalysis.GetLatestAnalysis) *AnalysisHandler {
	return &AnalysisHandler{run: run, latest: latest}
}

func (h *AnalysisHandler) Run(c *gin.Context) {
	report, err := h.run.Execute(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, portfolio.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, dto.NewAnalysisResponse(report))
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
	c.JSON(http.StatusOK, dto.NewAnalysisResponse(report))
}

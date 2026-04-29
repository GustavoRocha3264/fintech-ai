package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	apsnapshot "github.com/fintech/cbpi/backend-go/internal/application/snapshot"
	"github.com/fintech/cbpi/backend-go/internal/domain/snapshot"
)

type SnapshotHandler struct {
	history *apsnapshot.GetHistory
}

func NewSnapshotHandler(h *apsnapshot.GetHistory) *SnapshotHandler {
	return &SnapshotHandler{history: h}
}

type snapshotResponse struct {
	ID            string    `json:"id"`
	PortfolioID   string    `json:"portfolioId"`
	Timestamp     time.Time `json:"timestamp"`
	TotalValueBRL float64   `json:"totalValueBRL"`
	TotalValueUSD float64   `json:"totalValueUSD"`
}

func (h *SnapshotHandler) History(c *gin.Context) {
	items, err := h.history.Execute(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	out := make([]snapshotResponse, 0, len(items))
	for _, s := range items {
		out = append(out, toSnapshotResponse(s))
	}
	c.JSON(http.StatusOK, out)
}

func toSnapshotResponse(s snapshot.PortfolioSnapshot) snapshotResponse {
	return snapshotResponse{
		ID:            s.ID,
		PortfolioID:   s.PortfolioID,
		Timestamp:     s.Timestamp,
		TotalValueBRL: s.TotalValueBRL,
		TotalValueUSD: s.TotalValueUSD,
	}
}

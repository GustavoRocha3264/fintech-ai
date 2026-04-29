package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/fintech/cbpi/backend-go/internal/domain/portfolio"
)

type FXHandler struct {
	fx portfolio.FXRateProvider
}

func NewFXHandler(fx portfolio.FXRateProvider) *FXHandler {
	return &FXHandler{fx: fx}
}

type fxResponse struct {
	From      string    `json:"from"`
	To        string    `json:"to"`
	Rate      float64   `json:"rate"`
	FetchedAt time.Time `json:"fetchedAt"`
}

func (h *FXHandler) Get(c *gin.Context) {
	from := strings.ToUpper(c.Param("from"))
	to := strings.ToUpper(c.Param("to"))
	if from == "" || to == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from and to are required"})
		return
	}

	rate, err := h.fx.GetRate(from, to)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, fxResponse{
		From:      from,
		To:        to,
		Rate:      rate,
		FetchedAt: time.Now().UTC(),
	})
}

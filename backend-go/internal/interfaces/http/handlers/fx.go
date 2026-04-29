package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/fintech/cbpi/backend-go/internal/domain/portfolio"
	"github.com/fintech/cbpi/backend-go/internal/interfaces/http/dto"
)

type FXHandler struct {
	fx portfolio.FXRateProvider
}

func NewFXHandler(fx portfolio.FXRateProvider) *FXHandler {
	return &FXHandler{fx: fx}
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
	c.JSON(http.StatusOK, dto.NewFXResponse(from, to, rate))
}

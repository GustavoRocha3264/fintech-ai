package handlers

import (
	"errors"
	"net/http"

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
	c.JSON(http.StatusOK, dto.NewDashboardResponse(d))
}

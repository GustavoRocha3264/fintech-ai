package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	apportfolio "github.com/fintech/cbpi/backend-go/internal/application/portfolio"
	"github.com/fintech/cbpi/backend-go/internal/domain/portfolio"
	"github.com/fintech/cbpi/backend-go/internal/interfaces/http/dto"
)

type PortfolioHandler struct {
	create *apportfolio.CreatePortfolio
	get    *apportfolio.GetPortfolio
}

func NewPortfolioHandler(create *apportfolio.CreatePortfolio, get *apportfolio.GetPortfolio) *PortfolioHandler {
	return &PortfolioHandler{create: create, get: get}
}

func (h *PortfolioHandler) Create(c *gin.Context) {
	var req dto.CreatePortfolioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p, err := h.create.Execute(req.BaseCurrency)
	if err != nil {
		if errors.Is(err, portfolio.ErrInvalidBaseCurrency) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, toResponse(p))
}

func (h *PortfolioHandler) Get(c *gin.Context) {
	p, err := h.get.Execute(c.Param("id"))
	if err != nil {
		if errors.Is(err, portfolio.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, toResponse(p))
}

func toResponse(p *portfolio.Portfolio) dto.PortfolioResponse {
	return dto.PortfolioResponse{
		ID:           p.ID,
		BaseCurrency: p.BaseCurrency,
		CreatedAt:    p.CreatedAt,
	}
}

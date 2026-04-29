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
	create    *apportfolio.CreatePortfolio
	get       *apportfolio.GetPortfolio
	addPos    *apportfolio.AddPosition
	getValued *apportfolio.GetPortfolioWithValuation
}

func NewPortfolioHandler(
	create *apportfolio.CreatePortfolio,
	get *apportfolio.GetPortfolio,
	addPos *apportfolio.AddPosition,
	getValued *apportfolio.GetPortfolioWithValuation,
) *PortfolioHandler {
	return &PortfolioHandler{create: create, get: get, addPos: addPos, getValued: getValued}
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
	c.JSON(http.StatusCreated, dto.NewPortfolioResponse(p))
}

func (h *PortfolioHandler) Get(c *gin.Context) {
	p, err := h.get.Execute(c.Param("id"))
	if err != nil {
		writeRepoError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.NewPortfolioResponse(p))
}

func (h *PortfolioHandler) AddPosition(c *gin.Context) {
	var req dto.AddPositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pos, err := h.addPos.Execute(apportfolio.AddPositionInput{
		PortfolioID: c.Param("id"),
		Symbol:      req.Symbol,
		Quantity:    req.Quantity,
		Price:       req.Price,
		Currency:    req.Currency,
	})
	if err != nil {
		if errors.Is(err, portfolio.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, portfolio.ErrInvalidPositionCurrency) ||
			errors.Is(err, portfolio.ErrInvalidQuantity) ||
			errors.Is(err, portfolio.ErrInvalidPrice) ||
			errors.Is(err, portfolio.ErrEmptySymbol) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, dto.NewPositionResponse(*pos))
}

func (h *PortfolioHandler) GetWithValuation(c *gin.Context) {
	view, err := h.getValued.Execute(c.Param("id"))
	if err != nil {
		writeRepoError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.NewPortfolioWithValuationResponse(view))
}

func writeRepoError(c *gin.Context, err error) {
	if errors.Is(err, portfolio.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

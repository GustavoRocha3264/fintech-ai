package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/fintech/cbpi/backend-go/internal/interfaces/http/handlers"
)

func NewRouter(ph *handlers.PortfolioHandler, ah *handlers.AnalysisHandler) http.Handler {
	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())

	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	api := r.Group("/api/v1")
	{
		api.POST("/portfolios", ph.Create)
		api.GET("/portfolios/:id", ph.Get)
		api.POST("/portfolios/:id/positions", ph.AddPosition)
		api.GET("/portfolios/:id/valuation", ph.GetWithValuation)
		api.POST("/portfolios/:id/analysis", ah.Run)
		api.GET("/portfolios/:id/analysis/latest", ah.Latest)
	}

	return r
}

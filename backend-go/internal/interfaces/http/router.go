package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/fintech/cbpi/backend-go/internal/interfaces/http/handlers"
)

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func NewRouter(
	ph *handlers.PortfolioHandler,
	ah *handlers.AnalysisHandler,
	sh *handlers.SnapshotHandler,
	fh *handlers.FXHandler,
	dh *handlers.DashboardHandler,
) http.Handler {
	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger(), corsMiddleware())

	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	api := r.Group("/api/v1")
	{
		api.POST("/portfolios", ph.Create)
		api.GET("/portfolios/:id", ph.Get)
		api.POST("/portfolios/:id/positions", ph.AddPosition)
		api.GET("/portfolios/:id/valuation", ph.GetWithValuation)
		api.POST("/portfolios/:id/analysis", ah.Run)
		api.GET("/portfolios/:id/analysis/latest", ah.Latest)
		api.GET("/portfolios/:id/snapshots", sh.History)
		api.GET("/portfolios/:id/dashboard", dh.Get)
		api.GET("/fx/:from/:to", fh.Get)
	}

	return r
}

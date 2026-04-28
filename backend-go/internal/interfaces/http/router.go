package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/fintech/cbpi/backend-go/internal/interfaces/http/handlers"
)

func NewRouter(ph *handlers.PortfolioHandler) http.Handler {
	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())

	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	api := r.Group("/api/v1")
	{
		api.POST("/portfolios", ph.Create)
		api.GET("/portfolios/:id", ph.Get)
	}

	return r
}

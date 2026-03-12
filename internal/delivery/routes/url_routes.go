package routes

import (
	"github.com/gin-gonic/gin"
	handlers "github.com/mohamedkaram400/url-shortener/internal/delivery/handlers"
	"github.com/mohamedkaram400/url-shortener/internal/delivery/middlewares"
)

func RegisterUrlGenerationRoutes(rg *gin.RouterGroup, urlGenerationHandler *handlers.UrlGenerationHandler, secretKey string) {
	protected := rg.Group("/urls")
	protected.Use(middlewares.AuthMiddleware(secretKey))

	{
		protected.POST("/shorten", urlGenerationHandler.GenerateShortUrl)
		protected.POST("/analytics/:short_code", urlGenerationHandler.GenerateLinkAnalytics)
	}
}

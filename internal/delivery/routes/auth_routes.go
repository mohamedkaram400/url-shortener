package routes

import (
	"github.com/gin-gonic/gin"
	handlers "github.com/mohamedkaram400/url-shortener/internal/delivery/handlers"
	"github.com/mohamedkaram400/url-shortener/internal/delivery/middlewares"
	// middlewares "github.com/mohamedkaram400/url-shortener/internal/delivery/middlewares"
)

func RegisterUserAuthRoutes(rg *gin.RouterGroup, authHandler *handlers.AuthHandler, secretKey string) {
	auth := rg.Group("/auth")

	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
        auth.POST("/refresh-token/",  authHandler.RefreshToken)

		// 🔐 Protected routes
		protected := auth.Group("/")
		protected.Use(middlewares.AuthMiddleware(secretKey))
		{
			protected.POST("/logout", authHandler.Logout)
		}
	}
}

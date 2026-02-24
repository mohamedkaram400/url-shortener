package routes

import (
	"github.com/gin-gonic/gin"
	handlers "github.com/mohamedkaram400/url-shortener/internal/delivery/handlers"
	// middlewares "github.com/mohamedkaram400/url-shortener/internal/delivery/middlewares"
)

func RegisterUserAuthRoutes(rg *gin.RouterGroup, authHandler *handlers.AuthHandler) {
	auth := rg.Group("/auth")

	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
        auth.POST("/logout/:user_id", authHandler.Logout)
	}
}

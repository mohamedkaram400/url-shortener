package routes

import (
	"github.com/gin-gonic/gin"
	handlers "github.com/mohamedkaram400/url-shortener/internal/delivery/handlers"
)

func RegisterUserAuthRoutes(rg *gin.RouterGroup, authHandler *handlers.AuthHandler) {
	auth := rg.Group("/auth")

	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
        auth.POST("/refresh-token",  authHandler.RefreshToken)

		auth.POST("/logout", authHandler.Logout)
	}
}

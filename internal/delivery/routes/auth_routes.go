package auth

import (
	"github.com/gin-gonic/gin"
	AuthHandler "github.com/mohamedkaram400/url-shortener/internal/delivery/http"
	middlewares "github.com/mohamedkaram400/url-shortener/internal/delivery/middlewares"
)

func RegisterUserAuthRoutes(rg *gin.RouterGroup, authHandler *AuthHandler.AuthHandler) {
	auth := rg.Group("/auth")

	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
        auth.POST("/logout", middlewares.AdminJWTAuth(), authHandler.Logout)
	}
}

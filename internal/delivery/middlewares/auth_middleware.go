package middlewares

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mohamedkaram400/url-shortener/auth"
	"github.com/mohamedkaram400/url-shortener/pkg"
)

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		
		authHeader := c.GetHeader("Authorization")
		tokenString, err := pkg.ExtractTokenFromHeader(authHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		
		log.Println("RAW HEADER:", authHeader)
		log.Println("TOKEN STRING:", tokenString)
		log.Println("PARTS:", strings.Split(tokenString, "."))

		claims, err := auth.ValidateJWT(secret, tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if claims.TokenType != "access" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token type"})
			c.Abort()
			return
		}

		// Put userID inside context
		c.Set("userID", claims.UserID)

		c.Next()
	}
}
package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mohamedkaram400/url-shortener/internal/core/services"
)

func RateLimiterMiddleware(rateLimiterService *services.RateLimiterService) gin.HandlerFunc {
	return func (c *gin.Context) {

		ip := c.ClientIP()

		allowed, err := rateLimiterService.Allow(c.Request.Context(), ip, 100, time.Minute)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{
				"error": "internal error",
			})
			return
		}

		if !allowed {
			c.AbortWithStatusJSON(429, gin.H{
				"error": "too many requests",
			})
			return 
		}

		c.Next()
	}
}
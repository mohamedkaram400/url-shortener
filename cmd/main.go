package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mohamedkaram400/url-shortener/config"
)


func main() {

	// Load config
	config := config.LoadData()

	// Init router
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.Use(gin.Logger(), gin.Recovery())

	TestServer(router)

	StartServer(router, config)
}


func StartServer(router *gin.Engine, config *config.Config) {
	if err := router.Run(config.AppPort); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
	log.Println("🚀 App started on port", config.AppPort)
}

func TestServer(router *gin.Engine) {
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
}

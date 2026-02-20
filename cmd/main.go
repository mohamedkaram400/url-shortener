package main

import "github.com/gin-gonic/gin"

func TestServer(router *gin.Engine) {
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
}

func main() {

	// 6. Init router
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.Use(gin.Logger(), gin.Recovery())

	TestServer(router)

	router.Run(":4444")
}


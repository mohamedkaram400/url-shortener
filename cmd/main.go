package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mohamedkaram400/url-shortener/config"
	"github.com/mohamedkaram400/url-shortener/conn"
	"github.com/mohamedkaram400/url-shortener/db/seeders"
	"github.com/mohamedkaram400/url-shortener/internal/adapters/repositories"
	"github.com/mohamedkaram400/url-shortener/pkg"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/mohamedkaram400/url-shortener/internal/core/services"
	"github.com/mohamedkaram400/url-shortener/internal/delivery/handlers"
	"github.com/mohamedkaram400/url-shortener/internal/delivery/routes"
)


func main() {

	// Load config
	godotenv.Load()
	config := config.LoadData()

	// Connect DB, Redis Cache
	db, sqlDB := conn.ConnectPostgres(config.DatabaseURL)
	redisClient := conn.ConnectRedis(config.RedisServer)

	
	// Ensure the connection is closed when main() exits
	defer sqlDB.Close()
	defer redisClient.Close()

	// Run migrations
	pkg.RunMigrations(config.DatabaseURL)

	// Run seeders
	seeders.Run(db)

	// User Auth Module
	authUserRepo := repositories.NewAuthRepo(db)
	sessionRepo := repositories.NewSessionRepo(db)
	authUserService := services.NewAuthService(authUserRepo, sessionRepo, config.AccessTokenTime, config.RefrashTokenTime, config.JWTSecretKey)
	authUserHandler := handlers.NewAuthHandler(authUserService, config.BaseURL)

	// User Auth Module
	shortUrlRepo := repositories.NewUrlGenerationRepo(db)
	shortUrlService := services.NewUrlGenerationService(shortUrlRepo, config.ShortCodeLenght, config.BaseURL)
	urlGenerationHandler := handlers.NewUrlGenerationHandler(shortUrlService)


	// Init router
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.Use(gin.Logger(), gin.Recovery())
	api := router.Group("/api")

	// Public redirect route OUTSIDE api
	router.GET("/:code", urlGenerationHandler.Redirect)

	// Protected URL APIs
	routes.RegisterUserAuthRoutes(api, authUserHandler)
	routes.RegisterUrlGenerationRoutes(api, urlGenerationHandler, config.JWTSecretKey)


	TestServer(router)
	StartServer(router, config)
}


func StartServer(router *gin.Engine, config *config.Config) {
	// router.SetTrustedProxies([]string{"127.0.0.1"})
	if err := router.Run("0.0.0.0:" + config.AppPort); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
	log.Println("🚀 App started on port", config.AppPort)
}

func TestServer(router *gin.Engine) {
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
}


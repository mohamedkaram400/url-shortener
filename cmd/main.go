package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/mohamedkaram400/url-shortener/config"
	"github.com/mohamedkaram400/url-shortener/conn"
	"github.com/mohamedkaram400/url-shortener/db/seeders"
	"github.com/mohamedkaram400/url-shortener/internal/adapters/repositories"

	"github.com/mohamedkaram400/url-shortener/internal/core/services"
	"github.com/mohamedkaram400/url-shortener/internal/delivery/handlers"
	"github.com/mohamedkaram400/url-shortener/internal/delivery/routes"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)


func main() {

	// Load config
	config := config.LoadData()

	// Connect DB, Redis Cache
	db, sqlDB := conn.ConnectPostgres(config.DatabaseURL)

	redisClient := conn.ConnectRedis(config.RedisServer)

	
	// Ensure the connection is closed when main() exits
	defer sqlDB.Close()
	defer redisClient.Close()

	// Run migrations
	runMigrations(config.DatabaseURL)

	// Run seeders
	seeders.Run(db)

	// User Auth Module
	authUserRepo := repositories.NewAuthRepo(db)
	authUserService := services.NewAuthService(authUserRepo)
	authUserHandler := handlers.NewAuthHandler(authUserService)



	// Init router
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.Use(gin.Logger(), gin.Recovery())
	api := router.Group("/api")


	routes.RegisterUserAuthRoutes(api, authUserHandler)

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


func runMigrations(dbURL string) {
	m, err := migrate.New(
		"file://./db/migrations", 
		dbURL,
	)
	if err != nil {
		log.Fatal("❌ Failed to create migrate instance:", err)
	}

	err = m.Up() // runs all pending migrations
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal("❌ Failed to run migrations:", err)
	}

	log.Println("✅ Migrations ran successfully!")
}

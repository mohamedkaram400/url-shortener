package conn

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectPostgres(dbURL string) (*gorm.DB, *sql.DB) {

    // Open a connection to the database
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal("❌ Failed to connect database:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("❌ Failed to get underlying sql.DB:", err)
	}

	// Verify the connection is established and valid by pinging the database
	err = sqlDB.Ping()
    if err != nil {
        log.Fatal("❌ Error connecting to the database:", err)
    }

	log.Println("✅ Successfully connected to PostgreSQL!")

	return db, sqlDB
}
package conn

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectPostgres(dbURL string) (*gorm.DB, *sql.DB) {
	var db *gorm.DB
	var err error

	for i := 0; i < 10; i++ {
		// Open a connection to the database
		db, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})

		// Verify the connection is established and valid by pinging the database
		if err == nil {
			sqlDB, _ := db.DB()
			if err = sqlDB.Ping(); err == nil {
				log.Println("✅ Successfully connected to PostgreSQL!")
				return db, sqlDB
			}
		}

		log.Println("⏳ Waiting for database...")
		time.Sleep(2 * time.Second)
	}

	log.Fatal("❌ Failed to connect to database after retries:", err)
	return nil, nil
}
package seeders

import (
	"log"
	"gorm.io/gorm"
)

// Run all seeders here
func Run(db *gorm.DB) {
	log.Println("🚀 Running seeders...")

	SeedAdminUser(db)

	// SeedRoles(db)
	// SeedDefaultSettings(db)

	log.Println("✅ All seeders completed")
}
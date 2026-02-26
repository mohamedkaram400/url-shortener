package seeders

import (
	"errors"
	"log"

	"github.com/mohamedkaram400/url-shortener/internal/core/entities"
	"github.com/mohamedkaram400/url-shortener/pkg"
	"gorm.io/gorm"
)

func SeedAdminUser(db *gorm.DB) {
	// check if admin already exists
	var admin entities.User
	result := db.Where("username = ?", "Admin").Take(&admin)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Println("Creating Admin...")
		hashedPassword, err := pkg.HashPassword("Admin@123")
		if err != nil {
			log.Fatalf("❌ Error generating hashed password: %v", err)
		}

		admin = entities.User{
			UserName:   "Admin",
			FirstName:  "Ahmed",
			LastName:   "Ali",
			Email:    	"Admin@admin.com",
			Password:   hashedPassword,
		}

		if err := db.Create(&admin).Error; err != nil {
			log.Fatalf("❌ Failed to seed admin user: %v", err)
		}

		log.Println("✅ Admin user seeded successfully")
		return
	} else if result.Error != nil {
    	log.Fatalf("❌ Failed to check admin user: %v", result.Error)
	} else {
		log.Println("⚠️ Admin user already exists, skipping seeding")
	}
}
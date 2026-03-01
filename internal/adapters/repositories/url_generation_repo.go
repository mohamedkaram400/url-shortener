package repositories

import (
	"context"

	"github.com/mohamedkaram400/url-shortener/internal/core/entities"
	"gorm.io/gorm"
)


type UrlGenerationRepo struct {
	DB *gorm.DB
}

func NewUrlGenerationRepo(db *gorm.DB) (*UrlGenerationRepo) {
	return &UrlGenerationRepo{DB: db} 
}

func (r *UrlGenerationRepo) CreateUrl(ctx context.Context, url *entities.Url) error {
	return r.DB.WithContext(ctx).Create(url).Error
}
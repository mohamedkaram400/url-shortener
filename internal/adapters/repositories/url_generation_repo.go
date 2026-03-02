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

func (r *UrlGenerationRepo)	ShortCodeExists(ctx context.Context, shortCode string) (bool, error) {
    var count int64

    err := r.DB.WithContext(ctx).
        Model(&entities.Url{}).
        Where("short_code = ?", shortCode).
        Count(&count).Error

    if err != nil {
        return false, err
    }

    return count > 0, nil
}

func (r *UrlGenerationRepo) GetByShortCode(ctx context.Context, code string) (*entities.Url, error) {
	var url entities.Url

	err := r.DB.WithContext(ctx).
        Model(&entities.Url{}).
		Where("short_code = ?", code).
		First(&url).Error

	if err != nil {
        return nil, err
    }

	return  &url, nil
}
 
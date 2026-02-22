package repositories

import (
	"context"

	"github.com/mohamedkaram400/url-shortener/internal/core/entities"
	"gorm.io/gorm"
)

type AuthRepo struct {
	DB *gorm.DB
}

func newAuthRepo(db *gorm.DB) *AuthRepo {
	return &AuthRepo{DB: db}
}

func (r *AuthRepo) Register(ctx context.Context, user *entities.User) (*entities.User, error) {
	if err := r.DB.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	} 
	
	return user, nil
}

func (r *AuthRepo) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	if err := r.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AuthRepo) Logout(ctx context.Context, userID string) (string, error) {

	if err := r.DB.WithContext(ctx).Where("user_id = ?", userID).Delete(&entities.Session{}).Error; err != nil {
		return "", err
	}

	return "User logout success", nil
}

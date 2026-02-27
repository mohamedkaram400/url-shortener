package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/mohamedkaram400/url-shortener/internal/core/entities"
	"gorm.io/gorm"
)

type AuthRepo struct {
	DB *gorm.DB
}

func NewAuthRepo(db *gorm.DB) *AuthRepo {
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

func (r *AuthRepo) CheckEmailOrUsernameExists(
	ctx context.Context,
	email string,
	username string,
) (bool, bool, error) {

	var user entities.User

	err := r.DB.WithContext(ctx).
		Select("email", "username").
		Where("email = ? OR username = ?", email, username).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, false, nil
		}
		return false, false, err
	}

	emailExists := user.Email == email
	usernameExists := user.UserName == username

	return emailExists, usernameExists, nil
}

func (r *AuthRepo) CreateSession(ctx context.Context, session *entities.Session) error {
	return r.DB.WithContext(ctx).Create(session).Error
}

func (r *AuthRepo) DeleteSession(ctx context.Context, sessionID uint) (error) {
	if err := r.DB.WithContext(ctx).Where("id = ?", sessionID).Delete(&entities.Session{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *AuthRepo) GetValidSessionByRefreshToken(ctx context.Context, hashedToken string) (*entities.Session, error) {

	var session entities.Session

	err := r.DB.WithContext(ctx).
		Where("refresh_token = ?", hashedToken).
		Where("expires_at > ?", time.Now()).
		First(&session).Error

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *AuthRepo) Logout(ctx context.Context, HashedToken string) error {
	if err := r.DB.WithContext(ctx).Where("refresh_token = ?", HashedToken).Delete(&entities.Session{}).Error; err != nil {
		return err
	}

	return nil
}
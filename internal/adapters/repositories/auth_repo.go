package repositories

import (
	"context"
	"errors"

	"github.com/lib/pq"
	"github.com/mohamedkaram400/url-shortener/internal/core/entities"
	domainerrors "github.com/mohamedkaram400/url-shortener/internal/core/errors"
	"gorm.io/gorm"
)

type AuthRepo struct {
	DB *gorm.DB
}

func NewAuthRepo(db *gorm.DB) *AuthRepo {
	return &AuthRepo{DB: db}
}

func (r *AuthRepo) Register(ctx context.Context, user *entities.User) (*entities.User, error) {
	err := r.DB.WithContext(ctx).Create(user).Error

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" { // duplicate key
			return nil, domainerrors.ErrUserAlreadyExists
		}
		return nil, err
	}
	return user, nil
}

func (r *AuthRepo) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User

	err := r.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrUserNotFound
		}
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
            // No user exists → safe to register
            return false, false, nil
        }
		return false, false, err
	}

	emailExists := user.Email == email
	usernameExists := user.UserName == username

	return emailExists, usernameExists, nil
}

func (r *AuthRepo) Logout(ctx context.Context, HashedToken string) error {
	err := r.DB.WithContext(ctx).Where("refresh_token = ?", HashedToken).Delete(&entities.Session{}).Error
	if err != nil {
		return err
	}

	return nil
}
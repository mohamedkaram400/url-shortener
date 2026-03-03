package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/mohamedkaram400/url-shortener/internal/core/entities"
	domainerrors "github.com/mohamedkaram400/url-shortener/internal/core/errors"
	"gorm.io/gorm"
)

type SessionRepo struct {
	DB *gorm.DB
}

func NewSessionRepo(db *gorm.DB) *SessionRepo {
	return &SessionRepo{DB: db}
}

func (r *SessionRepo) CreateSession(ctx context.Context, session *entities.Session) error {
	err := r.DB.WithContext(ctx).Create(session).Error
	if err != nil {
		return err 
	}
	return nil
}

func (r *SessionRepo) DeleteSession(ctx context.Context, sessionID uint) (error) {
	err := r.DB.WithContext(ctx).Where("id = ?", sessionID).Delete(&entities.Session{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *SessionRepo) GetValidSessionByRefreshToken(ctx context.Context, hashedToken string) (*entities.Session, error) {

	var session entities.Session

	err := r.DB.WithContext(ctx).
		Where("refresh_token = ?", hashedToken).
		Where("expires_at > ?", time.Now()).
		First(&session).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrNotFound
		}
		return nil, err
	}

	return &session, nil
}
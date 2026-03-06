package ports

import (
	"context"

	"github.com/mohamedkaram400/url-shortener/internal/core/entities"
)

type SessionRepository interface {
	// Session
	CreateSession(ctx context.Context, session *entities.Session) (error)
	DeleteSession(ctx context.Context, session uint64) (error)
	GetValidSessionByRefreshToken(ctx context.Context, hashedToken string) (*entities.Session, error)
}

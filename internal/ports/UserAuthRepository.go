package ports

import (
	"context"

	"github.com/mohamedkaram400/url-shortener/internal/core/entities"
)

type UserAuthRepository interface {
	// Register
	Register(ctx context.Context, user *entities.User) (*entities.User, error)

	// Login
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
	CheckEmailOrUsernameExists(ctx context.Context, email string, username string) (bool, bool, error)

	// Session
	CreateSession(ctx context.Context, session *entities.Session) (error)
	DeleteSession(ctx context.Context, session uint) (error)
	GetValidSessionByRefreshToken(ctx context.Context, hashedToken string) (*entities.Session, error)

	// Logout
	Logout(ctx context.Context, refreshToken string) (error)
}

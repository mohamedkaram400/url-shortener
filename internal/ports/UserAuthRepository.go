package ports

import (
	"context"

	"github.com/mohamedkaram400/url-shortener/internal/core/entities"
)

type UserAuthRepository interface {
	Register(ctx context.Context, user *entities.User) (*entities.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
	CheckEmailOrUsernameExists(ctx context.Context, email string, username string) (bool, bool, error)
	CreateSession(ctx context.Context, session *entities.Session) (error)
	Logout(ctx context.Context, refreshToken string) (error)

}
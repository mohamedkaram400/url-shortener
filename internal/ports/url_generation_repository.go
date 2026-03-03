package ports

import (
	"context"

	"github.com/mohamedkaram400/url-shortener/internal/core/entities"
)

type UrlGenerationRepository interface {
	CreateUrl(ctx context.Context, url *entities.Url) (error)
	ShortCodeExists(ctx context.Context, shortCode string) (bool, error)
	GetByShortCode(ctx context.Context, code string) (*entities.Url, error)
	IncreaseCount(ctx context.Context, code string) (error)
}

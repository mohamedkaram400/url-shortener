package ports

import (
	"context"

	"github.com/mohamedkaram400/url-shortener/internal/core/entities"
)

type UrlGenerationRepository interface {
	CreateUrl(ctx context.Context, url *entities.Url) (error)
	ShortCodeExists(ctx context.Context, shortCode string) (bool, error)
	GetByShortCode(ctx context.Context, shortCode string) (*entities.Url, error)
	IncrementClickCount(ctx context.Context, urlID uint64) (error)
}

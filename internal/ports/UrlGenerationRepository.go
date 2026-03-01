package ports

import (
	"context"

	"github.com/mohamedkaram400/url-shortener/internal/core/entities"
)

type UrlGenerationRepository interface {
	CreateUrl(ctx context.Context, url *entities.Url) (error)
}

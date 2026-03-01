package services

import (
	"context"
	"time"

	"github.com/mohamedkaram400/url-shortener/internal/core/entities"
	"github.com/mohamedkaram400/url-shortener/internal/dto"
	"github.com/mohamedkaram400/url-shortener/internal/ports"
	"github.com/mohamedkaram400/url-shortener/pkg"
)

type UrlGenerationService struct {
	Repo ports.UrlGenerationRepository
	ShortCodeLenght int
}

func NewUrlGenerationService(repo ports.UrlGenerationRepository, shortCodeLenght int) *UrlGenerationService {
	return &UrlGenerationService{Repo: repo, ShortCodeLenght: shortCodeLenght}
}

func (s *UrlGenerationService) GenerateShortUrl(ctx context.Context, req *dto.ShortenUrlRequest) (*entities.Url, error) {

	var expiresAt *time.Time

	// only calculate if user provided expiration
	if req.ExpirationDays > 0 {
		t := time.Now().Add(time.Duration(req.ExpirationDays * 24 * uint(time.Hour)))
		expiresAt = &t
	}

	// Generate the short code
	shortCode, err := pkg.GenerateShortCode(s.ShortCodeLenght)
	if err != nil {
		return nil, err
	}

	url := entities.Url{
		OriginalURL:       req.LongUrl,
		Status:       	   req.Status,
		UserID:       	   req.UserId,
		CustomAlias:	   req.CustomAlias,
		ExpiresAt:    	   expiresAt,

		ShortCode:         shortCode,
	}

	if err := s.Repo.CreateUrl(ctx, &url);  err != nil {
		return nil, err
	}

	return &url, nil
}

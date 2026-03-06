package services

import (
	"context"
	"fmt"
	"time"

	"github.com/mohamedkaram400/url-shortener/internal/core/entities"
	domainerrors "github.com/mohamedkaram400/url-shortener/internal/core/errors"
	"github.com/mohamedkaram400/url-shortener/internal/dto"
	"github.com/mohamedkaram400/url-shortener/internal/ports"
	"github.com/mohamedkaram400/url-shortener/pkg"
)

type UrlGenerationService struct {
	Repo            ports.UrlGenerationRepository
	ShortCodeLenght int
	BaseURL         string
}

func NewUrlGenerationService(repo ports.UrlGenerationRepository, shortCodeLenght int, baseUrl string) *UrlGenerationService {
	return &UrlGenerationService{Repo: repo, ShortCodeLenght: shortCodeLenght, BaseURL: baseUrl}
}

func (s *UrlGenerationService) GenerateShortUrl(ctx context.Context, userID uint64, req *dto.ShortenUrlRequest) (string, error) {

	var expiresAt *time.Time
	var shortCode string

	// only calculate if user provided expiration
	if req.ExpirationDays > 0 {
		t := time.Now().AddDate(0, 0, req.ExpirationDays)
		expiresAt = &t
	}

	if req.CustomAlias != nil && *req.CustomAlias != "" {

		alias := *req.CustomAlias

		exists, err := s.Repo.ShortCodeExists(ctx, alias)
		if err != nil {
			return "", err
		}

		if exists {
			return "", domainerrors.ErrShortCodeExists
		}

		shortCode = alias
	} else {

		for {

			// Generate the short code
			code, err := pkg.GenerateShortCode(s.ShortCodeLenght)
			if err != nil {
				return "", err
			}

			exists, err := s.Repo.ShortCodeExists(ctx, code)
			if err != nil {
				return "", err
			}

			if !exists {
				shortCode = code
				break
			}
		}
	}

	url := entities.Url{
		OriginalURL: req.LongUrl,
		Status:      req.Status,
		UserID:      userID,
		CustomAlias: req.CustomAlias,

		ExpiresAt: expiresAt,
		ShortCode: shortCode,
	}

	if err := s.Repo.CreateUrl(ctx, &url); err != nil {
		return "", err
	}

	// Add custom alias which is the base URL and concatenate with short code generated to build the valid short URL
	shortUrl := fmt.Sprintf("%s/%s", s.BaseURL, shortCode)

	return shortUrl, nil
}

func (s *UrlGenerationService) Redirect(c context.Context, code string) (string, error) {
	url, err := s.Repo.GetByShortCode(c, code)
	if err != nil {
		return "", err
	}

	if url.Status != "Active" {
		return "", domainerrors.ErrURLInactive
	}

	if url.ExpiresAt != nil && time.Now().After(*url.ExpiresAt) {
		return "", domainerrors.ErrLinkExpired
	}

	if err := s.Repo.IncrementClickCount(c, url.ID); err != nil {
		return "", err
	}

	return url.OriginalURL, nil
}

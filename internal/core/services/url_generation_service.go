package services

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/mohamedkaram400/url-shortener/internal/core/entities"
	domainerrors "github.com/mohamedkaram400/url-shortener/internal/core/errors"
	"github.com/mohamedkaram400/url-shortener/internal/dto"
	"github.com/mohamedkaram400/url-shortener/internal/ports"
	"github.com/mohamedkaram400/url-shortener/internal/responses"
	"github.com/mohamedkaram400/url-shortener/pkg"
	"github.com/redis/go-redis/v9"
)

type UrlGenerationService struct {
	Repo            ports.UrlGenerationRepository
	Redis           *redis.Client
	ShortCodeLenght int
	BaseURL         string
}

func NewUrlGenerationService(repo ports.UrlGenerationRepository, redis *redis.Client, shortCodeLenght int, baseUrl string) *UrlGenerationService {
	return &UrlGenerationService{Repo: repo, Redis: redis, ShortCodeLenght: shortCodeLenght, BaseURL: baseUrl}
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

func (s *UrlGenerationService) Redirect(ctx context.Context, code string) (string, error) {
    var finalURL string
    var urlID uint64

	// 1. Check Redis cache
	url, err := s.getURLFromRedis(ctx, code)
	if err != nil {
		return "", err
	}

	if err == nil && url != "" {
        parts := strings.SplitN(url, "|", 2)

		id, _ := strconv.ParseUint(parts[0], 10, 64)

		finalURL = parts[1]
		urlID = id

		log.Println("long url, id from cache: ", finalURL, urlID)

    } else {

		// 2. Cache miss → query DB
		url, err := s.Repo.GetByShortCode(ctx, code)
		if err != nil {
			return "", err
		}

		if url.Status != "Active" {
			return "", domainerrors.ErrURLInactive
		}

		if url.ExpiresAt != nil && time.Now().After(*url.ExpiresAt) {
			return "", domainerrors.ErrLinkExpired
		}

		finalURL = url.OriginalURL
		urlID = url.ID

		log.Println("long url from db: ", finalURL)

		// 3. Store in Redis (non-blocking)
		_ = s.storeURLInRedis(ctx, code, finalURL, urlID)
	}
	
	// 4. Async analytics increment
	go func() {
		if urlID != 0 {
			if err := s.Repo.IncrementClickCount(context.Background(), urlID); err != nil {
				log.Println("click count update failed:", err)
			}
		}	
	}()

	return finalURL, nil
}

func (s *UrlGenerationService) storeURLInRedis(ctx context.Context, code string, longURL string, urlID uint64) error {
    value := fmt.Sprintf("%d|%s", urlID, longURL)
    return s.Redis.Set(ctx, "url:"+code, value, 24*time.Hour).Err()
}

func (s *UrlGenerationService) getURLFromRedis(ctx context.Context, code string) (string, error) {
	url, err := s.Redis.Get(ctx, "url:"+code).Result()

	if err == redis.Nil {
		return "", nil // cache miss
	}

	if err != nil {
		return "", err
	}

	return url, nil
}

func (s *UrlGenerationService) GenerateLinkAnalytics(ctx context.Context, code string) (*responses.LinkAnalyticsResponse, error) {
	var expires string

	url, err := s.Repo.GetByShortCode(ctx, code)
	if err != nil {
		return  nil, err
	}
	
	if url.ExpiresAt != nil {
		expires = url.ExpiresAt.Format("Monday, 02 January 2006 at 15:04")
	}

	return &responses.LinkAnalyticsResponse{
		LongURL:     url.OriginalURL,
		ShortURL:    url.ShortCode,
		TotalClicks: url.ClickCount,
		CustomAlias: url.CustomAlias,
		Status:      url.Status,
		CreatedAt:   expires,
		ExpiresAt:   url.ExpiresAt.Format("Monday, 02 January 2006 at 15:04"),
	}, nil
}

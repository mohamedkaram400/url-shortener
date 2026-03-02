package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	coreErrors "github.com/mohamedkaram400/url-shortener/internal/core/errors"
	"github.com/mohamedkaram400/url-shortener/internal/core/services"
	"github.com/mohamedkaram400/url-shortener/internal/dto"
)

type UrlGenerationHandler struct {
	UrlGenerationService *services.UrlGenerationService
}

func NewUrlGenerationHandler(urlGenerationService *services.UrlGenerationService) *UrlGenerationHandler {
	return &UrlGenerationHandler{UrlGenerationService: urlGenerationService}
}

func (h *UrlGenerationHandler) GenerateShortUrl(c *gin.Context) {
	// Get the request 
	var req dto.ShortenUrlRequest

	// Validate the request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": coreErrors.ValidationError(err)})
		return 
	}

	// Pass the request the generation short url service
	shortUrl, err := h.UrlGenerationService.GenerateShortUrl(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the response with short url
	c.JSON(http.StatusCreated, gin.H{"short_url": shortUrl, "long_url": req.LongUrl})
}

func (h *UrlGenerationHandler) Redirect(c *gin.Context) {
	code := c.Param("code")

	originalURL, err := h.UrlGenerationService.GetByShortCode(c, code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusMovedPermanently, originalURL) // 301
}
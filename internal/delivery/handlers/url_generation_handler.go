package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohamedkaram400/url-shortener/internal/core/services"
	"github.com/mohamedkaram400/url-shortener/internal/dto"
	"github.com/mohamedkaram400/url-shortener/internal/responses"
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
		respondError(c, err)
		return 
	}

	userID, _ := c.Get("userID")

	// Pass the request the generation short url service
	shortURL, err := h.UrlGenerationService.GenerateShortUrl(c.Request.Context(), userID.(uint64), &req)
	if err != nil {
		respondError(c, err)
		return
	}

	data := responses.URLResponse{
		ShortURL: shortURL,
		LongURL:  req.LongUrl,
	}

	// Return the response with short url
	c.JSON(http.StatusCreated, responses.SuccessResponse("URL Generated successfully", data))
}

func (h *UrlGenerationHandler) Redirect(c *gin.Context) {
	code := c.Param("code")

	originalURL, err := h.UrlGenerationService.Redirect(c.Request.Context(), code)
	if err != nil {
		respondError(c, err)
		return
	}

	c.Redirect(http.StatusFound, originalURL) // 302
}

func (h *UrlGenerationHandler) GenerateLinkAnalytics(c *gin.Context) {
	code := c.Param("short_code")

	analyticsData, err := h.UrlGenerationService.GenerateLinkAnalytics(c.Request.Context(), code)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse("URL Analytics Generated successfully", analyticsData))
}
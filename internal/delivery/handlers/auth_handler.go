package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	coreErrors "github.com/mohamedkaram400/url-shortener/internal/core/errors"
	"github.com/mohamedkaram400/url-shortener/internal/core/services"
	"github.com/mohamedkaram400/url-shortener/internal/dto"
	"github.com/mohamedkaram400/url-shortener/internal/responses"
)

type AuthHandler struct {
	AuthService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{AuthService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(422, gin.H{"error": coreErrors.ValidationError(err)})
		return 
	}

	user, accessToken, refreshToken, err := h.AuthService.Register(c, &req)
	if err != nil {

		if errors.Is(err, coreErrors.ErrEmailAlreadyExists) ||
			errors.Is(err, coreErrors.ErrUserNameAlreadyExists) {

			c.JSON(http.StatusConflict, gin.H{"error": err.Error(),})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": coreErrors.ValidationError(err),})
		return
	}

	c.JSON(201, responses.AuthResponse{
		Message: "User Register Successfully",
		AccessToken: accessToken,
		RefreshToken: refreshToken,
		User: responses.NewAuthResponse(user),
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusUnprocessableEntity, err)
		return 
	}

	ip := c.ClientIP()
	device := c.GetHeader("User-Agent")

	user, accessToken, refreshToken, err := h.AuthService.Login(c, &req, ip, device)
	if err != nil {
		switch {
		case errors.Is(err, coreErrors.ErrUserNotFound),
			errors.Is(err, coreErrors.ErrInvalidCredentials):
			respondError(c, http.StatusUnauthorized, err)
			return
		default:
			respondError(c, http.StatusInternalServerError, err)
			return
		}
	}

	c.JSON(200, responses.AuthResponse{
		Message: "User Login Successfully",
		AccessToken: accessToken,
		RefreshToken: refreshToken,
		User: responses.NewAuthResponse(user),
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
    authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		respondError(c, http.StatusUnauthorized, errors.New("invalid authorization header"))
		return
	}

	// Split "Bearer <token>"
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 && strings.ToLower(parts[0]) != "bearer" {
		respondError(c, http.StatusUnauthorized, errors.New("invalid authorization header"))
		return 
	}

	refreshToken := parts[1]

	message, err := h.AuthService.Logout(c, refreshToken)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return 	
	}

	c.JSON(http.StatusOK, gin.H{"message": message})
}

// =========================
// Helper: standard error response
// =========================
func respondError(c *gin.Context, status int, err error) {
	// Validation errors (DTO struct validation or custom FieldError)
	if ve := coreErrors.ValidationError(err); len(ve) > 0 {
		c.JSON(status, gin.H{"errors": ve})
		return
	}

	// fallback
	c.JSON(status, gin.H{"error": err.Error()})
}
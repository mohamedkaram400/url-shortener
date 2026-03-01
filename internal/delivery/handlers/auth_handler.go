package handlers

import (
	"errors"
	"net/http"

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

	ip := c.ClientIP()
	device := c.GetHeader("User-Agent")

	user, accessToken, refreshToken, maxAge, err := h.AuthService.Register(c, &req, ip, device)
	if err != nil {

		if errors.Is(err, coreErrors.ErrEmailAlreadyExists) ||
			errors.Is(err, coreErrors.ErrUserNameAlreadyExists) {

			c.JSON(http.StatusConflict, gin.H{"error": err.Error(),})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": coreErrors.ValidationError(err),})
		return
	}

	// Set refresh token in HttpOnly cookie
	c.SetCookie(
		"refresh_token",
		refreshToken,
		maxAge,
		"/",
		"url-shortener.test",
		false, // secure (true in HTTPS)
		true,  // HttpOnly
	)

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

	user, accessToken, refreshToken, maxAge, err := h.AuthService.Login(c, &req, ip, device)
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

	// Set refresh token in HttpOnly cookie
	c.SetCookie(
		"refresh_token",
		refreshToken,
		maxAge,
		"/",
		"url-shortener.test",
		false, // secure (true in HTTPS)
		true,  // HttpOnly
	)

	c.JSON(200, responses.AuthResponse{
		Message: "User Login Successfully",
		AccessToken: accessToken,
		RefreshToken: refreshToken,
		User: responses.NewAuthResponse(user),
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token missing"})
		return
	}

	message, err := h.AuthService.Logout(c, refreshToken)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return 	
	}

	// 🧹 Clear cookie
	c.SetCookie(
		"refresh_token",
		"",
		-1, // delete cookie
		"/",
		"url-shortener.test",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{"message": message})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Get refresh token from cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token missing"})
		return
	}

	ip := c.ClientIP()
	device := c.GetHeader("User-Agent")

	accessToken, refreshToken, maxAge, err := h.AuthService.RefreshToken(c, refreshToken, ip, device)
	if err != nil {
		respondError(c, http.StatusUnauthorized, err)
		return
	}

	// Set new rotated refresh token
	c.SetCookie(
		"refresh_token",
		refreshToken,
		maxAge,
		"/",
		"url-shortener.test",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"refresh_token": refreshToken,
	})
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
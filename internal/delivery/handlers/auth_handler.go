package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	coreErrors "github.com/mohamedkaram400/url-shortener/internal/core/errors"
	"github.com/mohamedkaram400/url-shortener/internal/core/services"
	"github.com/mohamedkaram400/url-shortener/internal/dto"
	"github.com/mohamedkaram400/url-shortener/internal/responses"
	"github.com/mohamedkaram400/url-shortener/pkg"
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
	refreshToken, err := pkg.ExtractTokenFromHeader(authHeader)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
		return 
	}

	fmt.Println("LOGOUT ERROR:", err)
	fmt.Printf("TYPE: %T\n", err)

	message, err := h.AuthService.Logout(c, refreshToken)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return 	
	}

	c.JSON(http.StatusOK, gin.H{"message": message})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken	string	`"json"refresh_token"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusUnprocessableEntity, err)
		return
	}

	ip := c.ClientIP()
	device := c.GetHeader("User-Agent")

	accessToken, refreshToken, err := h.AuthService.RefreshToken(c, req.RefreshToken, ip, device)
	if err != nil {
		respondError(c, http.StatusUnauthorized, err)
		return
	}

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
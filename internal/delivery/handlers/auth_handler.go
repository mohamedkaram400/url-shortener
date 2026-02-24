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

	user, accessToken, refreshToken, err := h.AuthService.Register(c, &req)
	if err != nil {

		if errors.Is(err, coreErrors.ErrEmailAlreadyExists) ||
			errors.Is(err, coreErrors.ErrUserNameAlreadyExists) {

			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": coreErrors.ValidationError(err),
		})
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
		c.JSON(422, gin.H{"error": coreErrors.ValidationError(err)})
		return 
	}

	ip := c.ClientIP()
	device := c.GetHeader("User-Agent")

	user, accessToken, refreshToken, err := h.AuthService.Login(c, &req, ip, device)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": coreErrors.ValidationError(err)})
		return 	
	}

	c.JSON(200, responses.AuthResponse{
		Message: "User Login Successfully",
		AccessToken: accessToken,
		RefreshToken: refreshToken,
		User: responses.NewAuthResponse(user),
	})
}


func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken := c.Param("refresh_token")

	message, err := h.AuthService.Logout(c, refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": coreErrors.ValidationError(err)})
		return 	
	}

	c.JSON(200, responses.AuthResponse{
		Message: message,
	})
}
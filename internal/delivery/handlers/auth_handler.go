package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	domainerrors "github.com/mohamedkaram400/url-shortener/internal/core/errors"
	"github.com/mohamedkaram400/url-shortener/internal/core/services"
	"github.com/mohamedkaram400/url-shortener/internal/dto"
	"github.com/mohamedkaram400/url-shortener/internal/responses"
)

type AuthHandler struct {
	AuthService *services.AuthService
	BaseURL		string
}

func NewAuthHandler(authService *services.AuthService, baseURL string) *AuthHandler {
	return &AuthHandler{AuthService: authService, BaseURL: baseURL}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, err)
		return 
	}

	user, accessToken, refreshToken, maxAge, err := h.AuthService.Register(c, &req, c.ClientIP(), c.GetHeader("User-Agent"))
	if err != nil {
		respondError(c, err)
		return
	}

	// Set refresh token in HttpOnly cookie
	c.SetCookie(
		"refresh_token",
		refreshToken,
		maxAge,
		"/",
		h.BaseURL,
		false, // secure (true in HTTPS)
		true,  // HttpOnly
	)

	data := responses.AuthUserResponse{
		Tokens: responses.TokensResponse(accessToken, refreshToken),
		User: responses.ToUserResponse(user),
	}

	c.JSON(http.StatusCreated, responses.SuccessResponse("User Register Successfully", data))
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, err)
		return 
	}

	user, accessToken, refreshToken, maxAge, err := h.AuthService.Login(c, &req, c.ClientIP(), c.GetHeader("User-Agent"))
	if err != nil {
        respondError(c, err)
        return
    }

	// Set refresh token in HttpOnly cookie
	c.SetCookie(
		"refresh_token",
		refreshToken,
		maxAge,
		"/",
		h.BaseURL,
		false, // secure (true in HTTPS)
		true,  // HttpOnly
	)

	data := responses.AuthUserResponse{
		Tokens: responses.TokensResponse(accessToken, refreshToken),
		User: responses.ToUserResponse(user),
	}

	c.JSON(http.StatusOK, responses.SuccessResponse("User Login Successfully", data))
}

func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token missing"})
		return
	}

	message, err := h.AuthService.Logout(c, refreshToken)
	if err != nil {
		respondError(c, err)
		return 	
	}

	// 🧹 Clear cookie
	c.SetCookie(
		"refresh_token",
		"",
		-1, // delete cookie
		"/",
		h.BaseURL,
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{"message": message})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Get refresh token from cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		respondError(c, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": ""})
		return
	}

	ip := c.ClientIP()
	device := c.GetHeader("User-Agent")

	accessToken, refreshToken, maxAge, err := h.AuthService.RefreshToken(c, refreshToken, ip, device)
	if err != nil {
		respondError(c, err)
		return
	}

	// Set new rotated refresh token
	c.SetCookie(
		"refresh_token",
		refreshToken,
		maxAge,
		"/",
		h.BaseURL,
		false,
		true,
	)

	data := responses.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	c.JSON(http.StatusOK, responses.SuccessResponse("Access Token Generated Successfully", data))
}


// =========================
// Helper: standard error response
// =========================
func respondError(c *gin.Context, err error) {
    // 1. Validation errors
    if ve := domainerrors.ValidationError(err); len(ve) > 0 {
        c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": ve})
        return
    }

    // 2. Domain errors mapping
    switch {
    case errors.Is(err, domainerrors.ErrUserNotFound),
         errors.Is(err, domainerrors.ErrNotFound):
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return

    case errors.Is(err, domainerrors.ErrInvalidCredentials):
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return

    case errors.Is(err, domainerrors.ErrUserAlreadyExists),
         errors.Is(err, domainerrors.ErrEmailAlreadyExists),
         errors.Is(err, domainerrors.ErrUserNameAlreadyExists),
         errors.Is(err, domainerrors.ErrRefreshTokenMissing),
         errors.Is(err, domainerrors.ErrShortCodeExists):
        c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
        return

    case errors.Is(err, domainerrors.ErrLinkExpired):
        c.JSON(http.StatusGone, gin.H{"error": err.Error()})
        return

    case errors.Is(err, domainerrors.ErrURLInactive):
        c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
        return
    }

    // 3. Fallback for unexpected/internal errors
    log.Println("Internal error:", err.Error()) // log internally
    c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
}
package services

import (
	"context"
	"time"

	"github.com/mohamedkaram400/url-shortener/auth"
	"github.com/mohamedkaram400/url-shortener/internal/core/entities"
	coreErrors "github.com/mohamedkaram400/url-shortener/internal/core/errors"
	"github.com/mohamedkaram400/url-shortener/internal/dto"
	"github.com/mohamedkaram400/url-shortener/internal/ports"
	"github.com/mohamedkaram400/url-shortener/pkg"
)


type AuthService struct {
	Repo ports.UserAuthRepository
	AccessTokenTime int
	RefrashTokenTime int
	JWTSecret string
}

func NewAuthService(authRepo ports.UserAuthRepository, accessTokenTime int, refrashTokenTime int, secret string) *AuthService {
	return &AuthService{Repo: authRepo, JWTSecret: secret, AccessTokenTime: accessTokenTime, RefrashTokenTime: refrashTokenTime}
}

func (s *AuthService) Register(ctx context.Context, req *dto.RegisterRequest) (*entities.User, string, string, error) {
	var fieldErrors []*coreErrors.FieldError

	emailExists, usernameExists, err := s.Repo.CheckEmailOrUsernameExists(ctx, req.Email, req.UserName)
	if err != nil {
		return nil, "", "", err
	}

	if emailExists {
		fieldErrors = append(fieldErrors, coreErrors.NewFieldError("email", "email already exists"))
	}
	if usernameExists {
		fieldErrors = append(fieldErrors, coreErrors.NewFieldError("username", "username already exists"))
	}

	if len(fieldErrors) > 0 {
		return nil, "", "", &coreErrors.MultiFieldError{Errors: fieldErrors}
	}

	hashedPassword, err := pkg.HashPassword(req.Password)
		if err != nil {
		return nil, "", "", err
	}

	user := &entities.User{
		UserName:  req.UserName,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  hashedPassword,
	}

	user, err = s.Repo.Register(ctx, user)
	if err != nil {
		return nil, "", "", err
	}

	accessToken, refreshToken, err := auth.IsuueTokens(s.JWTSecret, s.AccessTokenTime, s.RefrashTokenTime, user.ID)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest, ip string, device string) (*entities.User, string, string, error) {

	user, err := s.Repo.GetUserByEmail(ctx, req.Email)
	if err != nil || user == nil {
		return nil, "", "", coreErrors.ErrUserNotFound
	}

	// check password
	if err := pkg.CheckPassword(user.Password, req.Password); err != nil {
		return nil, "", "", coreErrors.ErrInvalidCredentials
	}

	accessToken, refreshToken, err := auth.IsuueTokens(s.JWTSecret, s.AccessTokenTime, s.RefrashTokenTime, user.ID)
	if err != nil {
		return nil, "", "", err
	}

	hashedToken := pkg.HashToken(refreshToken)

	session := entities.Session{
		UserID:       user.ID,
		Device:       device,
		RefreshToken: hashedToken,
		IpAddress:    ip,
		ExpiresAt:    time.Now().Add(time.Duration(s.RefrashTokenTime) * 24 * time.Hour),
	}

	if err := s.Repo.CreateSession(ctx, &session); err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) (string, error) {
	err := s.Repo.Logout(ctx, refreshToken)
	if err != nil {
		return "", err
	}

	return "User logout successfully", nil
}

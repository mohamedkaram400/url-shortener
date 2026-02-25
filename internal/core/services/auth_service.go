package services

import (
	"context"
	"time"

	"github.com/mohamedkaram400/url-shortener/internal/core/entities"
	coreErrors "github.com/mohamedkaram400/url-shortener/internal/core/errors"
	"github.com/mohamedkaram400/url-shortener/internal/dto"
	"github.com/mohamedkaram400/url-shortener/internal/ports"
	"github.com/mohamedkaram400/url-shortener/pkg"
)


type AuthService struct {
	Repo ports.UserAuthRepository
}

func NewAuthService(authRepo ports.UserAuthRepository) *AuthService {
	return &AuthService{Repo: authRepo}
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

	hashedPassword := pkg.HashPassword(req.Password)

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

	return user, "access", "refresh", nil
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

	accessToken := "generate-access-token"
	refreshToken := "generate-refresh-token"

	session := entities.Session{
		UserID:       user.ID,
		Device:       device,
		RefreshToken: refreshToken,
		IpAddress:    ip,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
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

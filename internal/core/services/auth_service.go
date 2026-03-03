package services

import (
	"context"
	"time"

	"github.com/mohamedkaram400/url-shortener/auth"
	"github.com/mohamedkaram400/url-shortener/internal/core/entities"
	coreErrors "github.com/mohamedkaram400/url-shortener/internal/core/errors"
	domainerrors "github.com/mohamedkaram400/url-shortener/internal/core/errors"
	"github.com/mohamedkaram400/url-shortener/internal/dto"
	"github.com/mohamedkaram400/url-shortener/internal/ports"
	"github.com/mohamedkaram400/url-shortener/pkg"
)


type AuthService struct {
	AuthRepo ports.UserAuthRepository
	SessionRepo ports.SessionRepository
	AccessTokenTime int
	RefrashTokenTime int
	JWTSecret string
}

func NewAuthService(authRepo ports.UserAuthRepository, sessionRepo ports.SessionRepository, accessTokenTime int, refrashTokenTime int, secretKey string) *AuthService {
	return &AuthService{AuthRepo: authRepo, SessionRepo: sessionRepo, JWTSecret: secretKey, AccessTokenTime: accessTokenTime, RefrashTokenTime: refrashTokenTime}
}

func (s *AuthService) Register(ctx context.Context, req *dto.RegisterRequest, ip string, device string) (*entities.User, string, string, int, error) {
	var fieldErrors []*coreErrors.FieldError

	emailExists, usernameExists, err := s.AuthRepo.CheckEmailOrUsernameExists(ctx, req.Email, req.UserName)
	if err != nil {
		return nil, "", "", 0, err
	}

	if emailExists {
		fieldErrors = append(fieldErrors, coreErrors.NewFieldError("email", "email already exists"))
	}
	if usernameExists {
		fieldErrors = append(fieldErrors, coreErrors.NewFieldError("username", "username already exists"))
	}

	if len(fieldErrors) > 0 {
		return nil, "", "", 0, &coreErrors.MultiFieldError{Errors: fieldErrors}
	}

	hashedPassword, err := pkg.HashPassword(req.Password)
		if err != nil {
		return nil, "", "", 0, err
	}

	user := &entities.User{
		UserName:  req.UserName,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  hashedPassword,
	}

	user, err = s.AuthRepo.Register(ctx, user)
	if err != nil {
		return nil, "", "", 0, err
	}

	accessToken, refreshToken, err := auth.IsuueTokens(s.JWTSecret, s.AccessTokenTime, s.RefrashTokenTime, user.ID)
	if err != nil {
		return nil, "", "", 0, err
	}

	hashedToken := pkg.HashToken(refreshToken)

	session := entities.Session{
		UserID:       user.ID,
		Device:       device,
		RefreshToken: hashedToken,
		IpAddress:    ip,
		ExpiresAt:    time.Now().Add(time.Duration(s.RefrashTokenTime) * 24 * time.Hour),
	}

	if err := s.SessionRepo.CreateSession(ctx, &session); err != nil {
		return nil, "", "", 0, err
	}

	maxAge := int(time.Duration(s.RefrashTokenTime) * 24 * time.Hour / time.Second)

	return user, accessToken, refreshToken, maxAge, nil
}

func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest, ip string, device string) (*entities.User, string, string, int, error) {

	user, err := s.AuthRepo.GetUserByEmail(ctx, req.Email)
	if err != nil || user == nil {
		return nil, "", "", 0, coreErrors.ErrNotFound
	}

	// check password
	if err := pkg.CheckPassword(user.Password, req.Password); err != nil {
		return nil, "", "", 0, coreErrors.ErrInvalidCredentials
	}

	accessToken, refreshToken, err := auth.IsuueTokens(s.JWTSecret, s.AccessTokenTime, s.RefrashTokenTime, user.ID)
	if err != nil {
		return nil, "", "", 0, err
	}

	hashedToken := pkg.HashToken(refreshToken)

	session := entities.Session{
		UserID:       user.ID,
		Device:       device,
		RefreshToken: hashedToken,
		IpAddress:    ip,
		ExpiresAt:    time.Now().Add(time.Duration(s.RefrashTokenTime) * 24 * time.Hour),
	}

	if err := s.SessionRepo.CreateSession(ctx, &session); err != nil {
		return nil, "", "", 0, err
	}

	maxAge := int(time.Duration(s.RefrashTokenTime) * 24 * time.Hour / time.Second)

	return user, accessToken, refreshToken, maxAge, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string, ip string, device string) (string, string, int, error) {

	// 1. Validate JWT first
	claims, err := auth.ValidateJWT(s.JWTSecret, refreshToken)
	if err != nil {
		return "", "", 0, err 
	}
 
	if claims.TokenType != "refresh" {
		return "", "", 0, domainerrors.ErrInvalidCredentials
	}

	// 2. Hash refresh token
	hashedToken := pkg.HashToken(refreshToken)

	// 3. Check session in DB
	session, err := s.SessionRepo.GetValidSessionByRefreshToken(ctx, hashedToken)
	if err != nil {
		return "", "", 0, err
	}

	// 4. ROTATION (delete old refresh token)
	if err := s.SessionRepo.DeleteSession(ctx, session.ID); err != nil {
		return "", "", 0, err
	}

	// 5. Issue new tokens
	newAccess, newRefresh, err := auth.IsuueTokens(
		s.JWTSecret,
		s.AccessTokenTime,
		s.RefrashTokenTime,
		session.UserID,
	)
	if err != nil {
		return "", "", 0, err
	}

	// 6. Store new hashed token
	newHashed := pkg.HashToken(newRefresh)

	newSession := entities.Session{
		UserID:       session.UserID,
		Device:       device,
		IpAddress:    ip,
		RefreshToken: newHashed,
		ExpiresAt:    time.Now().Add(time.Duration(s.RefrashTokenTime) * 24 * time.Hour),
	}
	maxAge := int(time.Duration(s.RefrashTokenTime) * 24 * time.Hour / time.Second)

	if err := s.SessionRepo.CreateSession(ctx, &newSession); err != nil {
		return "", "", 0, err
	}

	return newAccess, newRefresh, maxAge, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) (string, error) {
	claims, err := auth.ValidateJWT(s.JWTSecret, refreshToken)
	if err != nil {
		return "", err
	}

	if claims.TokenType != "refresh" {
		return "", domainerrors.ErrInvalidCredentials
	}

	hashedToken := pkg.HashToken(refreshToken)

	if err := s.AuthRepo.Logout(ctx, hashedToken); err != nil {
		return "", err
	}

	return "User logout successfully", nil
}

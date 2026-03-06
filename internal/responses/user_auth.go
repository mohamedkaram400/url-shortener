package responses

import (
	"time"

	"github.com/mohamedkaram400/url-shortener/internal/core/entities"
)

type AuthUserResponse struct {
	Tokens       TokenResponse       `json:"tokens"`
	User         UserResponse      	   `json:"user"`
}

type UserResponse struct {
	ID        uint64      `json:"id"`
	UserName  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func TokensResponse(accessToken string, refreshToken string) TokenResponse {
	return TokenResponse{
		AccessToken:        accessToken,
		RefreshToken:       refreshToken,
	}
}

func ToUserResponse(user *entities.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		UserName:  user.UserName,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}

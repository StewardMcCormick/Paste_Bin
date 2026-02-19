package domain

import (
	"github.com/StewardMcCormick/Paste_Bin/internal/dto"
	"time"
)

type APIKey struct {
	Key       string
	Prefix    string
	ExpiresAt time.Time
}

type User struct {
	Id        int64
	Username  string
	Password  string
	APIKey    APIKey
	CreatedAt time.Time
}

func (u *User) ToResponse() *dto.UserResponse {
	return &dto.UserResponse{
		Id:       u.Id,
		Username: u.Username,
		APIKey: dto.APIKeyResponse{
			Key:       u.APIKey.Key,
			ExpiresAt: u.APIKey.ExpiresAt,
		},
		CreatedAt: u.CreatedAt,
	}
}

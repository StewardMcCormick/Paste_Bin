package domain

import (
	"time"

	"github.com/StewardMcCormick/Paste_Bin/internal/dto"
)

type APIKey struct {
	Key       string
	Prefix    string
	CreatedAt time.Time
	ExpiresAt time.Time
}

func (k *APIKey) ToResponse() dto.APIKeyResponse {
	return dto.APIKeyResponse{
		Key:       k.Key,
		ExpiresAt: k.ExpiresAt,
	}
}

type User struct {
	Id        int64
	Username  string
	Password  string
	APIKey    APIKey
	CreatedAt time.Time
}

func (u *User) ToResponse() *dto.UserResponse {
	user := &dto.UserResponse{
		Id:        u.Id,
		Username:  u.Username,
		CreatedAt: u.CreatedAt,
		APIKey:    u.APIKey.ToResponse(),
	}

	return user
}

package dto

import "time"

type APIKeyResponse struct {
	Key       string    `json:"key"`
	ExpiresAt time.Time `json:"expires_at"`
}

type UserResponse struct {
	Id        int64          `json:"id"`
	Username  string         `json:"username"`
	APIKey    APIKeyResponse `json:"api_key"`
	CreatedAt time.Time      `json:"created_at"`
}

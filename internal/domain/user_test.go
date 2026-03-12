package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUser_ToResponse(t *testing.T) {
	value := User{
		Id:       10,
		Username: "User",
		Password: "pass",
		APIKey: APIKey{
			UserId:    10,
			Key:       "key",
			Prefix:    "prefix",
			CreatedAt: time.Now(),
			ExpiresAt: time.Now().Add(time.Hour),
		},
		CreatedAt: time.Now(),
	}

	response := value.ToResponse()

	assert.Equal(t, value.Id, response.Id)
	assert.Equal(t, value.Username, response.Username)
	assert.Equal(t, value.CreatedAt, response.CreatedAt)
	assert.Equal(t, value.APIKey.Key, response.APIKey.Key)
	assert.Equal(t, value.APIKey.ExpiresAt, response.APIKey.ExpiresAt)
}

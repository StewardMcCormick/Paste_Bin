package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPaste_ToResponse(t *testing.T) {
	value := &Paste{
		Id:           12,
		UserId:       10,
		Hash:         "paste_hash",
		Views:        15,
		Privacy:      ProtectedPolicy,
		PasswordHash: "password_hash",
		CreatedAt:    time.Now(),
		ExpireAt:     time.Now().Add(5 * time.Hour),
		Content:      "content",
	}

	response := value.ToResponse()

	assert.Equal(t, value.Id, response.Id)
	assert.Equal(t, value.Views, response.Views)
	assert.Equal(t, string(value.Privacy), response.Privacy)
	assert.Equal(t, value.CreatedAt, response.CreatedAt)
	assert.Equal(t, value.ExpireAt, response.ExpireAt)
	assert.Equal(t, value.Hash, response.Hash)
	assert.Equal(t, string(value.Content), response.Content)
}

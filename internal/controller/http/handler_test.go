package http

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	router = Router()
)

func TestHandler_HelloHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	body, _ := io.ReadAll(w.Body)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.Equal(t, `{"message": "Hello world!"}`, string(body))
}

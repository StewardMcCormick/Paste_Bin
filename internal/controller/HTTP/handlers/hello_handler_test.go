package handlers

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	handler = NewHandler()
)

func TestHelloHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handler.HelloHandler(w, req)

	body, _ := io.ReadAll(w.Body)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.Equal(t,
		`{"message":"Hello world!"}`+"\n",
		string(body),
	)
}

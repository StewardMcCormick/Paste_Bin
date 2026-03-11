package middleware

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONValidation_Handler_NoError(t *testing.T) {
	validJson := []byte(`{"message": "hello"}`)

	callCount := 0
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
		w.Write(validJson)
	})

	req := httptest.NewRequest("GET", "/", bytes.NewReader(validJson))
	w := httptest.NewRecorder()

	handler := NewJSONValidation()
	midd := handler.Handler(testHandler)

	midd.ServeHTTP(w, req)

	resultBody, err := io.ReadAll(w.Result().Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.Equal(t, string(validJson), string(resultBody))
	assert.Equal(t, 1, callCount)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}

type brokenBody struct {
}

func (b *brokenBody) Read(p []byte) (n int, err error) {
	return 0, errors.New("broken reader")
}

func TestJSONValidation_Handler_Error(t *testing.T) {
	cases := []struct {
		name         string
		value        io.Reader
		expectedCode int
	}{
		{
			"With broken body",
			&brokenBody{},
			http.StatusBadRequest,
		},
		{
			"Empty JSON",
			bytes.NewReader([]byte("")),
			http.StatusBadRequest,
		},
		{
			"Invalid JSON",
			bytes.NewReader([]byte(`{"message": ,}`)),
			http.StatusBadRequest,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			callCount := 0
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				callCount++
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest("GET", "/", tc.value)
			w := httptest.NewRecorder()

			handler := NewJSONValidation()
			midd := handler.Handler(testHandler)

			midd.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedCode, w.Result().StatusCode)
			assert.Equal(t, 0, callCount)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		})
	}
}

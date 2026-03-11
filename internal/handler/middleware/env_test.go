package middleware

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	cfgutil "github.com/StewardMcCormick/Paste_Bin/config/cfg_util"
	appctx "github.com/StewardMcCormick/Paste_Bin/internal/util/app_context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvironmental_Handler(t *testing.T) {
	cases := []struct {
		name     string
		value    string
		expected string
	}{
		{
			"With Prod env",
			string(cfgutil.ProductionEnv),
			string(cfgutil.ProductionEnv),
		},
		{
			"With Dev env",
			string(cfgutil.DevelopmentEnv),
			string(cfgutil.DevelopmentEnv),
		},
		{
			"With incorrect env",
			"incorrect env",
			string(cfgutil.ProductionEnv),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			callCount := 0
			body := []byte("Hello")
			var capturedContext context.Context
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				callCount++
				capturedContext = r.Context()
				w.WriteHeader(http.StatusOK)
				w.Write(body)
			})

			req := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()

			handler := NewEnv(cfgutil.Env(tc.value))
			midd := handler.Handler(testHandler)

			midd.ServeHTTP(w, req)

			resultBody, err := io.ReadAll(w.Result().Body)
			require.NoError(t, err)

			resultEnv, err := appctx.GetEnv(capturedContext)
			require.NoError(t, err)

			assert.Equal(t, http.StatusOK, w.Result().StatusCode)
			assert.Equal(t, string(body), string(resultBody))
			assert.Equal(t, tc.expected, string(resultEnv))
			assert.Equal(t, 1, callCount)
		})
	}
}

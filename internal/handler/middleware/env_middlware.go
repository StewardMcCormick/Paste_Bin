package middleware

import (
	"context"
	"github.com/StewardMcCormick/Paste_Bin/config/cfg_util"
	httpUtil "github.com/StewardMcCormick/Paste_Bin/internal/util/http_util"
	"net/http"
)

func EnvironmentalMiddleware(env cfgUtil.Env) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if env == "" || (env != cfgUtil.ProductionEnv && env != cfgUtil.DevelopmentEnv) {
				env = cfgUtil.ProductionEnv
			}

			ctx := context.WithValue(r.Context(), httpUtil.EnvKey, env)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

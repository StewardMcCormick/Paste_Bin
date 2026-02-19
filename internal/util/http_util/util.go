package http_util

import (
	"net/http"
)

type loggerCtxKey string
type requestIdCtxKey string

var (
	LoggerKey    loggerCtxKey    = "logger"
	RequestIdKey requestIdCtxKey = "request_id"
	EnvKey                       = "env"
)

type WriterWithStatusCode struct {
	http.ResponseWriter
	StatusCode int
}

func (wc *WriterWithStatusCode) WriteHeader(status int) {
	wc.StatusCode = status
	wc.ResponseWriter.WriteHeader(status)
}

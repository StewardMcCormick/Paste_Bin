package middleware

import (
	"net/http"
)

type writerWithStatusCode struct {
	http.ResponseWriter
	statusCode int
}

func (wc *writerWithStatusCode) WriteHeader(status int) {
	wc.statusCode = status
	wc.ResponseWriter.WriteHeader(status)
}

package main

import (
	"net/http"
	"strings"

	"github.com/sk409/goconst"
)

func allowCredentials(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(goconst.HTTP_HEADER_ACCESS_CONTROL_ALLOW_CREDENTIALS, "true")
		next.ServeHTTP(w, r)
	})
}

func allowHeaders(headers []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := strings.Join(headers, ",")
		w.Header().Set(goconst.HTTP_HEADER_ACCESS_CONTROL_ALLOW_HEADERS, h)
		next.ServeHTTP(w, r)
	})
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(goconst.HTTP_HEADER_ACCESS_CONTROL_ALLOW_ORIGIN, "http://localhost:6900")
		next.ServeHTTP(w, r)
	})
}

func preflight(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			return
		}
		next.ServeHTTP(w, r)
	})
}

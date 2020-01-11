package main

import (
	"net/http"

	"github.com/sk409/goconst"
)

func main() {
	headerContentType := []string{goconst.HTTP_HEADER_CONTENT_TYPE}
	http.Handle("/auth", cors(&authHandler{}))
	http.Handle("/login", cors(allowCredentials(allowHeaders(headerContentType, preflight(&loginHandler{})))))
	http.Handle("/logout", cors(allowCredentials(&logoutHandler{})))
	http.Handle("/register", cors(allowCredentials(allowHeaders(headerContentType, preflight(&registerHandler{})))))
	http.Handle("/user", cors(allowCredentials(&userHandler{})))
	http.Handle("/users", cors(&usersHandler{}))
	http.ListenAndServe(":3300", nil)
}

package main

import (
	"net/http"

	"github.com/sk409/goconst"
)

func main() {
	headerContentType := []string{goconst.HTTP_HEADER_CONTENT_TYPE}
	http.Handle("/auth", cors(allowCredentials(&authHandler{})))
	http.Handle("/downloads", cors(&downloadsHandler{}))
	http.Handle("/files", cors(allowMethods([]string{http.MethodGet, http.MethodPost, http.MethodPut}, preflight(&filesHandler{}))))
	http.Handle("/folders", cors(&foldersHandler{}))
	http.Handle("/follows/", cors(allowMethods([]string{http.MethodGet, http.MethodPost, http.MethodDelete}, preflight(&followsHandler{}))))
	http.Handle("/lessons/", cors(allowMethods([]string{http.MethodGet, http.MethodPost, http.MethodPut}, preflight(&lessonsHandler{}))))
	http.Handle("/login", cors(allowCredentials(allowHeaders(headerContentType, preflight(&loginHandler{})))))
	http.Handle("/materials/", cors(&materialsHandler{}))
	http.Handle("/logout", cors(allowCredentials(&logoutHandler{})))
	http.Handle("/register", cors(allowCredentials(allowHeaders(headerContentType, preflight(&registerHandler{})))))
	http.Handle("/user", cors(allowCredentials(&userHandler{})))
	http.Handle("/users/", cors(&usersHandler{}))
	http.Handle("/public/", http.FileServer(http.Dir("./")))
	http.ListenAndServe(":3300", nil)
}

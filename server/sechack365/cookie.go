package main

import "net/http"

func newCookie(name, value string, maxAge int) *http.Cookie {
	return &http.Cookie{
		Name:   name,
		Value:  value,
		MaxAge: maxAge,
		Path:   "/",
	}
}

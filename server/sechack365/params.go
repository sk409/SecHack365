package main

import (
	"os"
	"path/filepath"

	"github.com/sk409/gosession"
)

const (
	cookieMaxAge30Days  = 60 * 60 * 24 * 30
	cookieNameSessionID = "SESSION_ID"
	sessionKeyUser      = "SESSION_KEY_USER"
)

var (
	cwd            string
	pathLessons    string
	sessionManager = gosession.Manager{Provider: gosession.NewMemoryProvider()}
)

func init() {
	var err error
	cwd, err = os.Getwd()
	if err != nil {
		panic(err)
	}
	pathLessons = filepath.Join(cwd, "lessons")
}

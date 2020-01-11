package main

import (
	"github.com/sk409/gosession"
)

const (
	cookieMaxAge30Days  = 60 * 60 * 24 * 30
	cookieNameSessionID = "SESSION_ID"
	sessionKeyUser      = "SESSION_KEY_USER"
)

var (
	sessionManager = gosession.Manager{Provider: gosession.NewMemoryProvider()}
)

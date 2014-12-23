package main

import (
	"github.com/gorilla/sessions"
	"net/http"
)

type User struct {
	Username string
	Password string
}

// Authorization Key
var authKey = []byte("somesecret")

// Encryption Key
var encKey = []byte("someothersecret")

var store = sessions.NewCookieStore(authKey, encKey)

func initSession(r *http.Request) *sessions.Session {
	session, _ := store.Get(r, "auth")
	if session.IsNew {
		session.Options.Domain = "dead.coffee"
		session.Options.MaxAge = 0
		session.Options.HttpOnly = false
		session.Options.Secure = true
	}
	return session
}

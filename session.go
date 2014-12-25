package main

import (
	"github.com/gorilla/securecookie"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

// cookie handling

var hashKey = securecookie.GenerateRandomKey(64)
var blockKey = securecookie.GenerateRandomKey(32)
var cookieHandler = securecookie.New(hashKey, blockKey)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Session struct {
	User *User
}

func setSession(s *Session, w http.ResponseWriter) {
	if encoded, err := cookieHandler.Encode("session", s); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	} else {
		log.Println("Can't set cookie:", err.Error())
	}
}

func getSession(r *http.Request) *Session {
	if cookie, err := r.Cookie("session"); err == nil {
		s := &Session{}
		if err = cookieHandler.Decode("session", cookie.Value, s); err != nil {
			log.Println(err.Error())
			return nil
		}
		return s
	} else {
		log.Println(err.Error())
		return nil
	}
}

func clearSession(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		renderTemplate(w, r, "login_page", map[string]interface{}{})
	}
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")
		redirectTarget := "/login/"

		if b, err := bcrypt.GenerateFromPassword([]byte(password), 10); err == nil {
			log.Println("bcrypt hash:", string(b))
		}

		var user User
		if err := db.Read("users", username, &user); err == nil {
			if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err == nil {
				session := &Session{}
				session.User = &User{}
				session.User.Username = user.Username

				setSession(session, w)
				redirectTarget = "/index/"
			} else {
				log.Println("password doesn't match?", err.Error())
			}
		}
		log.Println("returning...")
		http.Redirect(w, r, redirectTarget, http.StatusFound)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	clearSession(w)
	http.Redirect(w, r, "/index/", 302)
}

package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/securecookie"
	"github.com/zenazn/goji/web"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
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

func setJwt() {
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	token.Claims["foo"] = "bar"
	token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	log.Println(token)
	//tokenString, err := token.SignedString(mySigningKey)
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

func loginHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	//if b, err := bcrypt.GenerateFromPassword([]byte(password), 10); err == nil {
	//	log.Println("bcrypt hash:", string(b))
	//}

	var user User
	if err := db.Read("users", username, &user); err == nil {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err == nil {
			session := &Session{}
			session.User = &User{}
			session.User.Username = user.Username
			setSession(session, w)
		} else {
			log.Println("password doesn't match?", err.Error())
		}
	}
	http.Redirect(w, r, r.Referer(), http.StatusFound)
}

func logoutHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	clearSession(w)
	http.Redirect(w, r, r.Referer(), http.StatusFound)
}

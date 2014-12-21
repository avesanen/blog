package main

//import "github.com/gorilla/sessions"

type User struct {
	Username string
	password string
}

type Session struct {
	User *User
}

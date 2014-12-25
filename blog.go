package main

import (
	"github.com/avesanen/fsdb"
	"log"
)

var db *fsdb.FsDb

func init() {
	f, err := fsdb.NewFsDb("./db")
	if err != nil {
		panic(err)
	}
	db = f
}

func main() {
	db.Write("users", "admin", &User{Username: "admin", Password: "$2a$10$76B/HUE7CTqMrsjrreMxNukqZw1VJBskAETxeKJ.SmsC9G9hmChpi"})
	log.Println("Blog started.")
	s := newServer("", "5000")
	s.Start()
}

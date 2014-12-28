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
	/*c := &Chapter{}
	c.Id = "123"
	c.Title = "Chaptertest"
	c.Pages = []string{"1", "2", "3", "4", "5", "6"}
	db.Write("chapters", "1", c)
	db.Write("chapters", "2", c)
	db.Write("chapters", "3", c)

	a := &Archive{}
	a.Id = "1"
	a.Title = "first chap"
	a.Chapters = []string{"1", "2", "3"}
	db.Write("archive", "1", &a)

	u := &User{}
	u.Username = "admin"
	u.Password = "$2a$10$76B/HUE7CTqMrsjrreMxNukqZw1VJBskAETxeKJ.SmsC9G9hmChpi"
	db.Write("users", "admin", &u)
	*/
	log.Println("Blog started.")
	s := newServer("", "5000")
	s.Start()
}

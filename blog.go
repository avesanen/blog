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
	log.Println("Blog started.")
	s := newServer("", "5000")
	s.Start()
}

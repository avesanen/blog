package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Server struct {
	host string
	port string
	r    *mux.Router
}

func NewServer(host, port string) *Server {
	s := &Server{
		r:    mux.NewRouter(),
		host: host,
		port: port,
	}
	return s
}

func (s *Server) Start() {
	s.r.StrictSlash(true)

	// Session
	s.r.Path("/login/").Methods("GET", "POST").HandlerFunc(loginHandler)

	// Articles
	s.r.Path("/{article}/edit").Methods("GET", "POST").HandlerFunc(articleEditHandler)
	s.r.Path("/{article}/").Methods("GET").HandlerFunc(articleViewHandler)
	s.r.Path("/").Methods("GET").HandlerFunc(indexHandler)

	// Files
	s.r.Path("/upload/").Methods("POST").HandlerFunc(uploadPostHandler)
	s.r.Path("/img/{file}").Methods("GET").HandlerFunc(viewImageHandler)
	s.r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	s.r.StrictSlash(true)
	http.Handle("/", s.r)
	http.ListenAndServe(s.host+":"+s.port, nil)
}

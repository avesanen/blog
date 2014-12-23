package main

import (
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"net/http"
)

// Server keeps up the http server
type Server struct {
	host string
	port string
}

func newServer(host string, port string) *Server {
	s := &Server{
		host: host,
		port: port,
	}
	return s
}

func (s *Server) Start() {
	// Root routes
	goji.Get("/", indexHandler)

	// Login routes
	goji.Get("/login/", loginHandler)
	goji.Post("/login/", loginHandler)

	// Admin console
	admin := web.New()
	goji.Handle("/admin/*", admin)

	// Image up/download
	goji.Post("/upload/", uploadPostHandler)
	goji.Get("/img/:file", viewImageHandler)

	// Article routes
	goji.Get("/:article/", articleViewHandler)
	goji.Get("/:article/edit", articleEditHandler)
	goji.Post("/:article/edit", articleEditHandler)

	// Static routes
	goji.Get("/*", http.FileServer(http.Dir("./static/")))

	// Serve static files
	goji.Serve()
}

/*
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
*/

package main

import (
	"flag"
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
	flag.Set("bind", host+":"+port)
	return s
}

func (s *Server) Start() {
	// Root routes
	goji.Get("/", indexHandler)

	// Login routes
	goji.Post("/u/login", loginHandler)
	goji.Get("/u/logout", requiresLogin(logoutHandler))

	// Admin console
	admin := web.New()
	goji.Handle("/admin/*", admin)

	// Image up/download
	goji.Post("/img/upload", requiresLogin(uploadPostHandler))
	goji.Get("/img/:file", viewImageHandler)

	// Article routes
	goji.Get("/a/:article", articleViewHandler)
	goji.Get("/a/:article/edit", requiresLogin(articleEditHandler))
	goji.Post("/a/:article/edit", requiresLogin(articleEditHandler))

	// Static routes
	goji.Get("/*", http.FileServer(http.Dir("./static/")))

	// Serve static files
	goji.Serve()
}

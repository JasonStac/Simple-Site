package server

import (
	"context"
	"fmt"
	"net/http"
)

type ErrorHandler func(w http.ResponseWriter, r *http.Request, err error)

type Server struct {
	port         int
	Router       *Router
	middlewares  []Middleware
	errorHandler ErrorHandler
	server       *http.Server
}

type ServerOption func(*Server)

func NewServer(port int) *Server {
	server := &Server{
		port:        port,
		middlewares: []Middleware{},
	}
	server.Router = newRouter(server)
	return server
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server.Router.serveHTTP(w, r)
}

func (server *Server) Run() error {
	addr := fmt.Sprintf(":%d", server.port)
	server.server = &http.Server{
		Addr:    addr,
		Handler: server,
	}

	fmt.Printf("Server starting on port %d\n", server.port)
	return server.server.ListenAndServe()
}

func (server *Server) Shutdown(ctx context.Context) error {
	return server.server.Shutdown(ctx)
}

func (server *Server) GetRouter() *Router {
	return server.Router
}

func WithErrorHandler(handler ErrorHandler) ServerOption {
	return func(server *Server) {
		server.errorHandler = handler
	}
}

func DefaultErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

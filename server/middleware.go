package server

import (
	"net/http"
)

type Middleware func(Handler) Handler

func chain(h Handler, middleware ...Middleware) Handler {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}
	return h
}

func (serv *Server) applyMiddleware(h Handler) Handler {
	if len(serv.middlewares) == 0 {
		return h
	}
	return chain(h, serv.middlewares...)
}

func (serv *Server) Use(middleware ...Middleware) {
	serv.middlewares = append(serv.middlewares, middleware...)
}

func LoggingMiddleware(next Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		println("Request:", r.Method, r.URL.Path)
		next(w, r)
		println("Response sent")
	}
}

func AuthMiddleware(next Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		//simple auth for now, update later
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

func RecoveryMiddleware(next Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next(w, r)
	}
}

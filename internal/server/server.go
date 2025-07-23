package server

import (
	"database/sql"
	"goserv/internal/handlers"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewServer(db *sql.DB, tmpl *template.Template) http.Handler {
	router := chi.NewRouter()

	// apply global middleware

	handlers.AddRoutes(router, db, tmpl)

	return router
}

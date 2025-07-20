package server

import (
	"database/sql"
	"goserv/internal/handlers"
	"html/template"
	"net/http"
)

func NewServer(db *sql.DB, tmpl *template.Template) http.Handler {
	mux := http.NewServeMux()
	handlers.AddRoutes(mux, db, tmpl)

	var handler http.Handler = mux

	// put global middleware on handler

	return handler
}

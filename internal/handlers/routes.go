package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func AddRoutes(r *chi.Mux, db *sql.DB, tmpl *template.Template) {
	userMiddleware := IsUser(db)

	r.Get("/", handleHome(db, tmpl))

	r.Get("/login", handleLogin(db, tmpl))
	r.Post("/login", handleLogin(db, tmpl))
	r.With(userMiddleware).Get("/logout", handleLogout(db, tmpl))
	r.With(userMiddleware).Post("/logout", handleLogout(db, tmpl))

	r.Route("/view", func(r chi.Router) {
		r.Get("/posts", handleViewContent(db, tmpl))
		r.Get("/tags", handleViewTags(db, tmpl))
		r.Get("/artists", handleViewArtists(db, tmpl))
	})

	r.With(userMiddleware).Route("/profile", func(r chi.Router) {
		r.Get("/", handleProfile(tmpl))
		r.Get("/post", handleAddContent(db, tmpl))
		r.Post("/post", handleAddContent(db, tmpl))
		r.Get("/uploads", handleViewUploads(db, tmpl))
		r.Get("/favourites", handleViewFavourites(db, tmpl))
	})

	r.Mount("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("styles"))))
	r.Mount("/assets/images/", http.StripPrefix("/assets/images/", http.FileServer(http.Dir("content"))))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("404 Not Found: %s\n", r.URL.Path)
		http.NotFound(w, r)
	})
}

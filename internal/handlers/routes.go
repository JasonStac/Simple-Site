package handlers

import (
	"database/sql"
	"html/template"
	"net/http"
)

func AddRoutes(mux *http.ServeMux, db *sql.DB, tmpl *template.Template) {
	mux.Handle("/view", handleViewContent(db, tmpl))
	mux.Handle("/register", handleRegister(db, tmpl))
	mux.Handle("/login", handleLogin(db, tmpl))
	mux.Handle("/logout", isUser(db, handleLogout(db, tmpl)))

	mux.Handle("/profile", isUser(db, handleProfile(tmpl)))
	mux.Handle("/add", isUser(db, handleAddContent(db, tmpl)))
	mux.Handle("/uploads", isUser(db, handleViewUploads(db, tmpl)))

	mux.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("./styles"))))
	mux.Handle("/assets/images/", http.StripPrefix("/assets/images/", http.FileServer(http.Dir("./content"))))

	mux.Handle("/", handleHome(db, tmpl))
}

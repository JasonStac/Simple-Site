package handlers

import (
	"database/sql"
	"goserv/utils"
	"net/http"
)

func isUser(db *sql.DB, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := utils.GetSessionUser(db, r)
		if user != nil {
			h.ServeHTTP(w, r)
			return
		}
		//display need login page
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})
}

func isAdmin(db *sql.DB, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := utils.GetSessionUser(db, r)
		if user != nil && user.IsAdmin {
			h.ServeHTTP(w, r)
			return
		}
		http.NotFound(w, r)
	})
}

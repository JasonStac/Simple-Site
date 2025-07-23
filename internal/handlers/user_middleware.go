package handlers

import (
	"context"
	"database/sql"
	"goserv/utils"
	"net/http"
)

func IsUser(db *sql.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := utils.GetSessionUser(db, r)
			if user != nil {
				ctx := context.WithValue(r.Context(), "user", user.Username)
				ctx = context.WithValue(ctx, "userID", user.ID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			//TODO: display need login
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		})
	}
}

func isUser(db *sql.DB, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user := utils.GetSessionUser(db, r); user != nil {
			h.ServeHTTP(w, r)
			return
		}
		//display need login page
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})
}

func isAdmin(db *sql.DB, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user := utils.GetSessionUser(db, r); user != nil && user.IsAdmin {
			h.ServeHTTP(w, r)
			return
		}
		http.NotFound(w, r)
	})
}

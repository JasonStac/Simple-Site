package handlers

import (
	"database/sql"
	"goserv/internal/models"
	"goserv/utils"
	"html/template"
	"net/http"
)

func handleHome(db *sql.DB, tmpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		user := utils.GetSessionUser(db, r)
		err := tmpl.ExecuteTemplate(w, "home.html", struct{ User *models.User }{
			User: user,
		})
		if err != nil {
			http.Error(w, "Template error", http.StatusInternalServerError)
		}
	})
}

package handlers

import (
	"database/sql"
	"goserv/internal/dao"
	"goserv/internal/models"
	"goserv/utils"
	"html/template"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func handleRegister(db *sql.DB, tmpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			user := utils.GetSessionUser(db, r)
			//TODO: add thingy if they're already logged in
			err := tmpl.ExecuteTemplate(w, "register.html", struct{ User *models.User }{
				User: user,
			})
			if err != nil {
				http.Error(w, "Template error", http.StatusInternalServerError)
			}
			return

		case http.MethodPost:
			username := r.FormValue("username")
			password := r.FormValue("password")

			hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, "Error hashing password", http.StatusInternalServerError)
				return
			}

			user := &models.User{Username: username, PassHash: string(hashedPass)}
			err = dao.AddUser(db, user)
			if err != nil {
				http.Error(w, "Error creating user", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/login", http.StatusSeeOther)

		default:
			http.Error(w, "Unsupported Method Request", http.StatusMethodNotAllowed)
			return
		}
	})
}

func handleLogin(db *sql.DB, tmpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			user := utils.GetSessionUser(db, r)
			//TODO: add thingy if they're already logged in
			err := tmpl.ExecuteTemplate(w, "login.html", struct{ User *models.User }{
				User: user,
			})
			if err != nil {
				http.Error(w, "Template error", http.StatusInternalServerError)
			}
			return

		case http.MethodPost:
			username := r.FormValue("username")
			password := r.FormValue("password")

			//TODO: Differentiate between user not existing and actual error occurring
			user, err := dao.GetUser(db, username)
			if err != nil {
				http.Error(w, "Error finding user", http.StatusInternalServerError)
				return
			}

			err = bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(password))
			if err != nil {
				http.Error(w, "Invalid username or password", http.StatusUnauthorized)
				return
			}

			sessionID := utils.GenerateSessionID()
			err = dao.SaveSession(db, username, sessionID)
			if err != nil {
				http.Error(w, "Error saving session", http.StatusInternalServerError)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "id",
				Value:    sessionID,
				Path:     "/",
				HttpOnly: true,
				Expires:  time.Now().Add(24 * time.Hour),
			})

			http.Redirect(w, r, "/", http.StatusSeeOther)

		default:
			http.Error(w, "Unsupported Method Request", http.StatusMethodNotAllowed)
			return
		}
	})
}

func handleLogout(db *sql.DB, tmpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			err := tmpl.ExecuteTemplate(w, "logout.html", nil)
			if err != nil {
				http.Error(w, "Template error", http.StatusInternalServerError)
				return
			}

		case http.MethodPost:
			confirmation := r.FormValue("yes")
			if confirmation == "" {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			cookie, err := r.Cookie("id")
			if err != nil {
				http.Error(w, "Failed to read session cookie", http.StatusBadRequest)
				return
			}

			username, err := dao.GetSessionUsername(db, cookie.Value)
			if err != nil {
				http.Error(w, "Failed to find session", http.StatusInternalServerError)
				return
			}

			err = dao.DeleteSession(db, username)
			if err != nil {
				http.Error(w, "Failed to delete session", http.StatusInternalServerError)
				return
			}

			cookie.MaxAge = -1
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return

		default:
			http.Error(w, "Unsupported Method Request", http.StatusMethodNotAllowed)
			return
		}
	})
}

func handleProfile(tmpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.ExecuteTemplate(w, "profile.html", nil)
		if err != nil {
			http.Error(w, "Template error", http.StatusInternalServerError)
			return
		}
	})
}

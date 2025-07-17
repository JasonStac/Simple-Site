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

type UserHandler struct {
	dao  *dao.UserDao
	tmpl *template.Template
}

func NewUserHandler(db *sql.DB, tmpl *template.Template) *UserHandler {
	return &UserHandler{dao: dao.NewUserDao(db), tmpl: tmpl}
}

func (h *UserHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		err := h.tmpl.ExecuteTemplate(w, "register.html", nil)
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
		err = h.dao.AddUser(user)
		if err != nil {
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)

	default:
		http.Error(w, "Unsupported Method Request", http.StatusMethodNotAllowed)
		return
	}

}

func (h *UserHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		err := h.tmpl.ExecuteTemplate(w, "login.html", nil)
		if err != nil {
			http.Error(w, "Template error", http.StatusInternalServerError)
		}
		return

	case http.MethodPost:
		username := r.FormValue("username")
		password := r.FormValue("password")

		//TODO: Fix differentiating between user not existing and actual error occurring
		user, err := h.dao.GetUser(username)
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
		err = h.dao.SaveSession(username, sessionID)
		if err != nil {
			http.Error(w, "Error saving session", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session",
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
}

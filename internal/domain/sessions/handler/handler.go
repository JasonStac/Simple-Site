package handler

import (
	"goserv/internal/domain/sessions/service"
	"goserv/internal/domain/users"
	"goserv/internal/middleware"
	"html/template"
	"net/http"
	"time"
)

type SessionHandler struct {
	svc  *service.SessionService
	tmpl *template.Template
}

func NewSessionHandler(svc *service.SessionService, tmpl *template.Template) *SessionHandler {
	return &SessionHandler{svc: svc, tmpl: tmpl}
}

func (h *SessionHandler) DisplayLogin(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if ok && userID != 0 {
		http.Error(w, "Already logged in", http.StatusUnauthorized)
		return
	}

	err := h.tmpl.ExecuteTemplate(w, "login.html", struct{ User *users.User }{
		User: &users.User{},
	})
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}

func (h *SessionHandler) Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	sessionID, err := h.svc.Login(r.Context(), username, password)
	if err != nil {
		http.Error(w, "Error logging in", http.StatusInternalServerError)
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
}

func (h *SessionHandler) DisplayLogout(w http.ResponseWriter, r *http.Request) {
	err := h.tmpl.ExecuteTemplate(w, "logout.html", struct{ User *users.User }{
		User: &users.User{},
	})
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}

func (h *SessionHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if confirmation := r.FormValue("yes"); confirmation == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	cookie, err := r.Cookie("id")
	if err != nil {
		http.Error(w, "Failed to read cookie", http.StatusUnauthorized)
		return
	}

	err = h.svc.Logout(r.Context(), cookie.Value)
	if err != nil {
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	cookie.MaxAge = -1
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

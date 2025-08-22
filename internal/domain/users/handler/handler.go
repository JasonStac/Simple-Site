package handler

import (
	"goserv/internal/domain/users"
	"goserv/internal/domain/users/service"
	"goserv/internal/middleware"
	"html/template"
	"net/http"
)

type UserHandler struct {
	svc  *service.UserService
	tmpl *template.Template
}

func NewUserHandler(svc *service.UserService, tmpl *template.Template) *UserHandler {
	return &UserHandler{svc: svc, tmpl: tmpl}
}

func (h *UserHandler) DisplayRegister(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if ok && userID != 0 {
		http.Error(w, "Already logged in", http.StatusUnauthorized)
	}

	err := h.tmpl.ExecuteTemplate(w, "register.html", struct{ User *users.User }{
		User: &users.User{},
	})
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	err := h.svc.Register(r.Context(), username, password)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	if err := h.tmpl.ExecuteTemplate(w, "profile.html", nil); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}

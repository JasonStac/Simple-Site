package handlers

import (
	"html/template"
	"net/http"
)

type HomeHandler struct {
	tmpl *template.Template
}

func NewHomeHandler(tmpl *template.Template) *HomeHandler {
	return &HomeHandler{tmpl: tmpl}
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.tmpl.ExecuteTemplate(w, "home.html", nil)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
	}
}

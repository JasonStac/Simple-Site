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

func (h *HomeHandler) serveHTTP(w http.ResponseWriter) {
	err := h.tmpl.ExecuteTemplate(w, "home.html", nil)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
	}
}

func (h *HomeHandler) HandleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	h.serveHTTP(w)
}

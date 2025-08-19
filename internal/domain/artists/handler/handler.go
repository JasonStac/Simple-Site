package handler

import (
	"goserv/internal/domain/artists/service"
	"goserv/internal/middleware"
	"html/template"
	"net/http"
)

type ArtistHandler struct {
	svc  *service.ArtistService
	tmpl *template.Template
}

func NewArtistHandler(svc *service.ArtistService, tmpl *template.Template) *ArtistHandler {
	return &ArtistHandler{svc: svc, tmpl: tmpl}
}

func (h *ArtistHandler) ListArtists(w http.ResponseWriter, r *http.Request) {
	artists, err := h.svc.ListArtists(r.Context())
	if err != nil {
		http.Error(w, "Failed to list tags", http.StatusInternalServerError)
		return
	}

	response := make([]string, len(artists))
	for i := range artists {
		response[i] = artists[i].Name
	}

	isUser := false
	username, ok := middleware.GetUserID(r)
	if ok && username != -1 {
		isUser = true
	}

	err = h.tmpl.ExecuteTemplate(w, "artists.html", struct {
		Artists []string
		IsUser  bool
	}{
		Artists: response,
		IsUser:  isUser,
	})
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}

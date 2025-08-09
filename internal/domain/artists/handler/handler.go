package handler

import (
	"goserv/internal/domain/artists/service"
	"goserv/internal/domain/users"
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

	var response []string
	for _, artist := range artists {
		name := artist.Name
		response = append(response, name)
	}

	var user *users.User
	username, ok := middleware.GetUserID(r)
	if !ok || username == -1 {
		user = nil
	} else {
		user = &users.User{}
	}

	err = h.tmpl.ExecuteTemplate(w, "artists.html", struct {
		Artists []string
		User    *users.User
	}{
		Artists: response,
		User:    user,
	})
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}

package handler

import (
	"goserv/internal/domain/tags/service"
	"goserv/internal/middleware"
	"html/template"
	"net/http"
)

type TagHandler struct {
	svc  *service.TagService
	tmpl *template.Template
}

func NewTagHandler(svc *service.TagService, tmpl *template.Template) *TagHandler {
	return &TagHandler{svc: svc, tmpl: tmpl}
}

func (h *TagHandler) ListTags(w http.ResponseWriter, r *http.Request) {
	tags, err := h.svc.ListTags(r.Context())
	if err != nil {
		http.Error(w, "Failed to list tags", http.StatusInternalServerError)
		return
	}

	response := make([]string, len(tags))
	for i := range tags {
		response[i] = tags[i].Name
	}

	isUser := false
	userID, ok := middleware.GetUserID(r)
	if ok && userID != 0 {
		isUser = true
	}

	err = h.tmpl.ExecuteTemplate(w, "tags.html", struct {
		Tags   []string
		IsUser bool
	}{
		Tags:   response,
		IsUser: isUser,
	})
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}

func (h *TagHandler) ListGeneralTags(w http.ResponseWriter, r *http.Request) {
	tags, err := h.svc.ListGeneralTags(r.Context())
	if err != nil {
		http.Error(w, "Failed to list tags", http.StatusInternalServerError)
		return
	}

	response := make([]string, len(tags))
	for i := range tags {
		response[i] = tags[i].Name
	}

	isUser := false
	userID, ok := middleware.GetUserID(r)
	if ok && userID != 0 {
		isUser = true
	}

	err = h.tmpl.ExecuteTemplate(w, "tags.html", struct {
		Tags   []string
		IsUser bool
	}{
		Tags:   response,
		IsUser: isUser,
	})
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}

func (h *TagHandler) ListPeopleTags(w http.ResponseWriter, r *http.Request) {
	tags, err := h.svc.ListPeopleTags(r.Context())
	if err != nil {
		http.Error(w, "Failed to list tags", http.StatusInternalServerError)
		return
	}

	response := make([]string, len(tags))
	for i := range tags {
		response[i] = tags[i].Name
	}

	isUser := false
	userID, ok := middleware.GetUserID(r)
	if ok && userID != 0 {
		isUser = true
	}

	err = h.tmpl.ExecuteTemplate(w, "people.html", struct {
		People []string
		IsUser bool
	}{
		People: response,
		IsUser: isUser,
	})
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}

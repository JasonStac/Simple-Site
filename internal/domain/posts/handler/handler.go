package handler

import (
	"errors"
	"goserv/internal/domain/posts"
	"goserv/internal/domain/posts/service"
	"goserv/internal/middleware"
	"goserv/internal/models"
	"goserv/internal/utils"
	myErrors "goserv/internal/utils/errors"
	"html/template"
	"net/http"
	"path"
)

type ResponseEntry struct {
	Path string
	ID   int
}

type PostHandler struct {
	svc  *service.PostService
	tmpl *template.Template
}

func NewPostHandler(svc *service.PostService, tmpl *template.Template) *PostHandler {
	return &PostHandler{svc: svc, tmpl: tmpl}
}

func (h *PostHandler) ViewAddPost(w http.ResponseWriter, r *http.Request) {
	err := h.tmpl.ExecuteTemplate(w, "add.html", struct{ MediaTypes []string }{
		MediaTypes: models.MediaType("").Values(),
	})
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}

func (h *PostHandler) AddPost(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Failed to authenticate user", http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	fileMedia := r.FormValue("media")
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	post := &posts.Post{Title: title, MediaType: models.MediaType(fileMedia), Filename: header.Filename}
	err = h.svc.AddPost(r.Context(), post, file, userID)
	if err != nil {
		http.Error(w, "Failed to add post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/profile/uploads", http.StatusSeeOther)
}

func (h *PostHandler) ListPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.svc.ListPosts(r.Context())
	if err != nil {
		http.Error(w, "Error listing posts", http.StatusInternalServerError)
		return
	}

	var responses []ResponseEntry
	for _, post := range posts {
		path := path.Join(post.Filename[0:2], post.Filename[2:4], post.Filename)
		responses = append(responses, ResponseEntry{Path: path, ID: post.ID})
	}

	isUser := false
	userID, ok := middleware.GetUserID(r)
	if ok && userID != -1 {
		isUser = true
	}

	err = h.tmpl.ExecuteTemplate(w, "list.html", struct {
		Posts  []ResponseEntry
		IsUser bool
	}{
		Posts:  responses,
		IsUser: isUser,
	})
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}

func (h *PostHandler) ViewPost(w http.ResponseWriter, r *http.Request) {
	postID, err := utils.GetPostIDFromPath(r.URL.Path)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	post, err := h.svc.GetPost(r.Context(), postID)
	if err != nil {
		if errors.Is(err, myErrors.ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Error getting post", http.StatusInternalServerError)
		return
	}

	path := path.Join(post.Filename[0:2], post.Filename[2:4], post.Filename)

	userID, ok := middleware.GetUserID(r)
	isUser := false
	if !ok && userID != -1 {
		isUser = true
	}

	err = h.tmpl.ExecuteTemplate(w, "view.html", struct {
		Path   string
		IsUser bool
	}{
		Path:   path,
		IsUser: isUser,
	})
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}

func (h *PostHandler) ListUserPosts(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}

	posts, err := h.svc.ListUserPosts(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to list posts", http.StatusInternalServerError)
		return
	}

	var responses []ResponseEntry
	for _, post := range posts {
		path := path.Join(post.Filename[0:2], post.Filename[2:4], post.Filename)
		responses = append(responses, ResponseEntry{Path: path, ID: post.ID})
	}

	err = h.tmpl.ExecuteTemplate(w, "uploads.html", struct{ Posts []ResponseEntry }{
		Posts: responses,
	})
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}

func (h *PostHandler) ListUserFavs(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}

	posts, err := h.svc.ListUserFavs(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to list posts", http.StatusInternalServerError)
		return
	}

	var responses []ResponseEntry
	for _, post := range posts {
		path := path.Join(post.Filename[0:2], post.Filename[2:4], post.Filename)
		responses = append(responses, ResponseEntry{Path: path, ID: post.ID})
	}

	err = h.tmpl.ExecuteTemplate(w, "uploads.html", struct{ Posts []ResponseEntry }{
		Posts: responses,
	})
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}

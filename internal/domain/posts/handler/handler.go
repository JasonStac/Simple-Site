package handler

import (
	"encoding/json"
	"errors"
	"goserv/internal/domain/artists"
	aService "goserv/internal/domain/artists/service"
	"goserv/internal/domain/posts"
	pService "goserv/internal/domain/posts/service"
	"goserv/internal/domain/tags"
	tService "goserv/internal/domain/tags/service"
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
	postSvc   *pService.PostService
	tagSvc    *tService.TagService
	artistSvc *aService.ArtistService
	tmpl      *template.Template
}

func NewPostHandler(
	postSvc *pService.PostService,
	tagSvc *tService.TagService,
	artistSvc *aService.ArtistService,
	tmpl *template.Template,
) *PostHandler {
	return &PostHandler{
		postSvc:   postSvc,
		tagSvc:    tagSvc,
		artistSvc: artistSvc,
		tmpl:      tmpl,
	}
}

func (h *PostHandler) ViewAddPost(w http.ResponseWriter, r *http.Request) {
	tagList, err := h.tagSvc.ListTags(r.Context())
	if err != nil {
		http.Error(w, "Error getting tags", http.StatusInternalServerError)
		// intentionally let continue for now
	}

	artistList, err := h.artistSvc.ListArtists(r.Context())
	if err != nil {
		http.Error(w, "Error gettings artists", http.StatusInternalServerError)
		// intentionally let continue for now
	}

	err = h.tmpl.ExecuteTemplate(w, "add.html", struct {
		MediaTypes []string
		TagList    []tags.Tag
		ArtistList []artists.Artist
	}{
		MediaTypes: models.MediaType("").Values(),
		TagList:    tagList,
		ArtistList: artistList,
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

	jsonTags := r.FormValue("tags")
	var tags []tags.Tag
	if err := json.Unmarshal([]byte(jsonTags), &tags); err != nil {
		http.Error(w, "Failed to read tags", http.StatusBadRequest)
		return
	}

	jsonArtists := r.FormValue("artists")
	var artists []artists.Artist
	if err := json.Unmarshal([]byte(jsonArtists), &artists); err != nil {
		http.Error(w, "Failed to read artists", http.StatusBadRequest)
		return
	}

	for _, t := range tags {
		if t.ID == 0 {
			id, err := h.tagSvc.AddTag(r.Context(), t.Name)
			if err != nil {
				http.Error(w, "Failed to add tag", http.StatusInternalServerError)
				return
			}
			t.ID = id
		}
	}

	for _, a := range artists {
		if a.ID == 0 {
			id, err := h.artistSvc.AddArtist(r.Context(), a.Name)
			if err != nil {
				http.Error(w, "Failed to add artist", http.StatusInternalServerError)
				return
			}
			a.ID = id
		}
	}

	post := &posts.Post{Title: title, MediaType: models.MediaType(fileMedia), Filename: header.Filename, Tags: &tags, Artists: &artists}
	err = h.postSvc.AddPost(r.Context(), post, file, userID)
	if err != nil {
		http.Error(w, "Failed to add post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/profile/uploads", http.StatusSeeOther)
}

func (h *PostHandler) ListPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.postSvc.ListPosts(r.Context())
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

	post, err := h.postSvc.GetPost(r.Context(), postID)
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
		Path    string
		IsUser  bool
		Artists []artists.Artist
		Tags    []tags.Tag
	}{
		Path:    path,
		IsUser:  isUser,
		Artists: *post.Artists,
		Tags:    *post.Tags,
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

	posts, err := h.postSvc.ListUserPosts(r.Context(), userID)
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

	posts, err := h.postSvc.ListUserFavs(r.Context(), userID)
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

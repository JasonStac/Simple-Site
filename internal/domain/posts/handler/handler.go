package handler

import (
	"errors"
	"goserv/internal/domain/posts"
	pService "goserv/internal/domain/posts/service"
	"goserv/internal/domain/tags"
	tService "goserv/internal/domain/tags/service"
	"goserv/internal/middleware"
	"goserv/internal/static/constant"
	"goserv/internal/static/enum"
	myErrors "goserv/internal/utils/errors"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

type ResponseEntry struct {
	Filename string
	FileExt  string
	ID       int
}

type PostHandler struct {
	postSvc *pService.PostService
	tagSvc  *tService.TagService
	tmpl    *template.Template
}

func NewPostHandler(
	postSvc *pService.PostService,
	tagSvc *tService.TagService,
	tmpl *template.Template,
) *PostHandler {
	return &PostHandler{
		postSvc: postSvc,
		tagSvc:  tagSvc,
		tmpl:    tmpl,
	}
}

func (h *PostHandler) ViewAddPost(w http.ResponseWriter, r *http.Request) {
	tagList, err := h.tagSvc.ListTags(r.Context())
	if err != nil {
		http.Error(w, "Error getting tags", http.StatusInternalServerError)
		// intentionally let continue for now
	}

	peopleList, err := h.tagSvc.ListPeopleTags(r.Context())
	if err != nil {
		http.Error(w, "Error getting people", http.StatusInternalServerError)
		// intentionally let continue for now
	}

	err = h.tmpl.ExecuteTemplate(w, "add.html", struct {
		MediaTypes []string
		GeneralTag string
		PeopleTag  string
		TagList    []tags.Tag
		PeopleList []tags.Tag
	}{
		MediaTypes: enum.MediaType("").Values(),
		GeneralTag: string(enum.TagGeneral),
		PeopleTag:  string(enum.TagPeople),
		TagList:    tagList,
		PeopleList: peopleList,
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

	tags, ok := middleware.GetTags(r)
	if !ok {
		http.Error(w, "Error reading tags", http.StatusInternalServerError)
		return
	}

	post := &posts.Post{Title: title, MediaType: enum.MediaType(fileMedia), Filename: header.Filename, Tags: tags}
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

	content := make([]ResponseEntry, len(posts))
	for i := range posts {
		content[i] = ResponseEntry{
			Filename: posts[i].Filename,
			FileExt:  constant.ThumbnailExt,
			ID:       posts[i].ID,
		}
	}

	isUser := false
	userID, ok := middleware.GetUserID(r)
	if ok && userID != 0 {
		isUser = true
	}

	err = h.tmpl.ExecuteTemplate(w, "list.html", struct {
		Posts  []ResponseEntry
		IsUser bool
	}{
		Posts:  content,
		IsUser: isUser,
	})
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}

func (h *PostHandler) ViewPost(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/view/posts/"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	isUser := false
	userID, ok := middleware.GetUserID(r)
	if ok && userID != 0 {
		isUser = true
	}

	post, isFav, err := h.postSvc.GetPostWithFavouriteStatus(r.Context(), postID, userID)
	if err != nil {
		if errors.Is(err, myErrors.ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Error getting post", http.StatusInternalServerError)
		return
	}

	tagMap, err := h.tagSvc.SeperateTagTypes(r.Context(), post.Tags)
	if err != nil {
		http.Error(w, "Error handling tags", http.StatusInternalServerError)
		return
	}

	err = h.tmpl.ExecuteTemplate(w, "view.html", struct {
		Filename  string
		FileExt   string
		ID        int
		IsUser    bool
		IsFav     bool
		People    []tags.Tag
		Tags      []tags.Tag
		Type      string
		TypeImage string
		TypeVideo string
	}{
		Filename:  post.Filename,
		FileExt:   post.FileExt[1:],
		ID:        postID,
		IsUser:    isUser,
		IsFav:     isFav,
		People:    tagMap[enum.TagPeople],
		Tags:      tagMap[enum.TagGeneral],
		Type:      string(post.MediaType),
		TypeImage: string(enum.MediaImage),
		TypeVideo: string(enum.MediaVideo),
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

	content := make([]ResponseEntry, len(posts))
	for i := range posts {
		content[i] = ResponseEntry{
			Filename: posts[i].Filename,
			FileExt:  constant.ThumbnailExt,
			ID:       posts[i].ID,
		}
	}

	err = h.tmpl.ExecuteTemplate(w, "uploads.html", struct{ Posts []ResponseEntry }{
		Posts: content,
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

	content := make([]ResponseEntry, len(posts))
	for i := range posts {
		content[i] = ResponseEntry{
			Filename: posts[i].Filename,
			FileExt:  constant.ThumbnailExt,
			ID:       posts[i].ID,
		}
	}

	err = h.tmpl.ExecuteTemplate(w, "favourites.html", struct{ Posts []ResponseEntry }{
		Posts: content,
	})
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	postID, ok := middleware.GetPostID(r)
	if !ok {
		http.Error(w, "Error reading post ID", http.StatusBadRequest)
		return
	}

	filename, fileExt, ok := middleware.GetCompleteFilename(r)
	if !ok {
		http.Error(w, "Error reading filename", http.StatusInternalServerError)
		return
	}

	err := h.postSvc.DeletePost(r.Context(), postID, filename, fileExt)
	if err != nil {
		http.Error(w, "Error deleting post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/profile/uploads", http.StatusSeeOther)
}

func (h *PostHandler) FavouritePost(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Error reading user id", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, "Error getting post id", http.StatusBadRequest)
		return
	}

	err = h.postSvc.FavouritePost(r.Context(), postID, userID)
	if err != nil {
		http.Error(w, "Error favouriting post", http.StatusInternalServerError)
		return
	}
}

func (h *PostHandler) UnfavouritePost(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "Error reading user id", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, "Error getting post id", http.StatusBadRequest)
		return
	}

	err = h.postSvc.UnfavouritePost(r.Context(), postID, userID)
	if err != nil {
		http.Error(w, "Error unfavouriting post", http.StatusInternalServerError)
		return
	}
}

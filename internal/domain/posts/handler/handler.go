package handler

import (
	"encoding/json"
	"errors"
	"goserv/internal/domain/posts"
	pService "goserv/internal/domain/posts/service"
	"goserv/internal/domain/tags"
	tService "goserv/internal/domain/tags/service"
	"goserv/internal/middleware"
	"goserv/internal/models"
	"goserv/internal/utils"
	myErrors "goserv/internal/utils/errors"
	"html/template"
	"log"
	"net/http"
	"path"
	"strconv"
)

type ResponseEntry struct {
	Path string
	ID   int
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
		MediaTypes: models.MediaType("").Values(),
		GeneralTag: string(models.TagGeneral),
		PeopleTag:  string(models.TagPeople),
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

	jsonGeneralTags := r.FormValue("tags")
	log.Printf("All general: %s\n\n", jsonGeneralTags)
	var generalTags []tags.Tag
	if jsonGeneralTags != "" {
		if err := json.Unmarshal([]byte(jsonGeneralTags), &generalTags); err != nil {
			http.Error(w, "Failed to read tags", http.StatusBadRequest)
			return
		}
	}

	jsonPeopleTags := r.FormValue("people")
	log.Printf("All people: %s\n\n", jsonPeopleTags)
	var peopleTags []tags.Tag
	if jsonPeopleTags != "" {
		if err := json.Unmarshal([]byte(jsonPeopleTags), &peopleTags); err != nil {
			http.Error(w, "Failed to read people", http.StatusBadRequest)
			return
		}
	}

	allTags := append(generalTags, peopleTags...)
	for i := range allTags {
		if allTags[i].ID == 0 {
			id, err := h.tagSvc.AddTag(r.Context(), allTags[i].Name, models.TagType(allTags[i].Type))
			if err != nil {
				http.Error(w, "Failed to add tag", http.StatusInternalServerError)
				return
			}
			allTags[i].ID = id
		}
	}

	post := &posts.Post{Title: title, MediaType: models.MediaType(fileMedia), Filename: header.Filename, Tags: allTags}
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

	paths := make([]ResponseEntry, len(posts))
	for i := range posts {
		paths[i] = ResponseEntry{
			Path: path.Join(posts[i].Filename[0:2], posts[i].Filename[2:4], posts[i].Filename),
			ID:   posts[i].ID,
		}
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
		Posts:  paths,
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

	isUser := false
	userID, ok := middleware.GetUserID(r)
	if ok && userID != -1 {
		isUser = true
	}

	post, isFav, err := h.postSvc.GetPostWithFavourite(r.Context(), postID, userID)
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

	path := path.Join(post.Filename[0:2], post.Filename[2:4], post.Filename)

	err = h.tmpl.ExecuteTemplate(w, "view.html", struct {
		Path   string
		ID     int
		IsUser bool
		IsFav  bool
		People []tags.Tag
		Tags   []tags.Tag
	}{
		Path:   path,
		ID:     postID,
		IsUser: isUser,
		IsFav:  isFav,
		People: tagMap[models.TagPeople],
		Tags:   tagMap[models.TagGeneral],
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

	paths := make([]ResponseEntry, len(posts))
	for i := range posts {
		paths[i] = ResponseEntry{
			Path: path.Join(posts[i].Filename[0:2], posts[i].Filename[2:4], posts[i].Filename),
			ID:   posts[i].ID,
		}
	}

	err = h.tmpl.ExecuteTemplate(w, "uploads.html", struct{ Posts []ResponseEntry }{
		Posts: paths,
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

	paths := make([]ResponseEntry, len(posts))
	for i := range posts {
		paths[i] = ResponseEntry{
			Path: path.Join(posts[i].Filename[0:2], posts[i].Filename[2:4], posts[i].Filename),
			ID:   posts[i].ID,
		}
	}

	err = h.tmpl.ExecuteTemplate(w, "favourites.html", struct{ Posts []ResponseEntry }{
		Posts: paths,
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

	filepath, ok := middleware.GetFilepath(r)
	if !ok {
		http.Error(w, "Error reading filepath", http.StatusInternalServerError)
		return
	}

	err := h.postSvc.DeletePost(r.Context(), postID, filepath)
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

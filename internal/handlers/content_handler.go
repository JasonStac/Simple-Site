package handlers

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"goserv/internal/dao"
	"goserv/internal/models"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type ContentHandler struct {
	dao  *dao.ContentDao
	tmpl *template.Template
}

func NewContentHandler(db *sql.DB, tmpl *template.Template) *ContentHandler {
	return &ContentHandler{dao: dao.NewContentDao(db), tmpl: tmpl}
}

// TODO: add encoding stuff, figure out GET vs POST request
func (h *ContentHandler) HandleAddContent(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.tmpl.ExecuteTemplate(w, "add.html", models.MediaTypes)

	case http.MethodPost:
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")
		artist := r.FormValue("artist")
		fileMedia, err := strconv.Atoi(r.FormValue("media"))
		if err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Failed to get uploaded file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		tempFile, err := os.CreateTemp("tmp", "upload-*")
		if err != nil {
			http.Error(w, "Failed to create temp file", http.StatusInternalServerError)
			return
		}
		defer os.Remove(tempFile.Name())
		defer tempFile.Close()

		hasher := sha256.New()
		multiWriter := io.MultiWriter(tempFile, hasher)

		_, err = io.Copy(multiWriter, file)
		tempFile.Close()
		if err != nil {
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			return
		}

		hashBytes := hasher.Sum(nil)
		hashHex := hex.EncodeToString(hashBytes)
		ext := strings.ToLower(filepath.Ext(header.Filename))

		dir1 := hashHex[0:2]
		dir2 := hashHex[2:4]
		finalDir := filepath.Join("content", dir1, dir2)
		finalName := hashHex + ext
		finalPath := filepath.Join(finalDir, finalName)

		if err := os.MkdirAll(finalDir, 0755); err != nil {
			http.Error(w, "Failed to create storage directory", http.StatusInternalServerError)
			return
		}

		content := &models.Content{Title: title, FileMedia: models.Media(fileMedia), Filename: finalName, Artist: artist}
		err = h.dao.AddContent(content)
		if err != nil {
			http.Error(w, "Failed to insert into DB", http.StatusInternalServerError)
			return
		}

		if err := os.Rename(tempFile.Name(), finalPath); err != nil {
			_ = h.dao.DeleteContentByFilename(finalName)
			http.Error(w, "Failed to store file", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/view", http.StatusSeeOther)
		return

	default:
		http.Error(w, "Unsupported Method Request", http.StatusMethodNotAllowed)
		return
	}
}

// TODO: add filtered viewing
func (h *ContentHandler) HandleViewContent(w http.ResponseWriter, r *http.Request) {
	content, err := h.dao.GetContentFiles()
	if err != nil {
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}

	h.tmpl.ExecuteTemplate(w, "view.html", struct{ Paths []string }{
		Paths: content,
	})
}

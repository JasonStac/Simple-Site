package handlers

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"goserv/internal/dao"
	"goserv/internal/models"
	"goserv/utils"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func handleAddContent(db *sql.DB, tmpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			err := tmpl.ExecuteTemplate(w, "add.html", struct{ MediaTypes []models.MediaType }{
				MediaTypes: models.MediaTypes,
			})
			if err != nil {
				http.Error(w, "Template error", http.StatusInternalServerError)
				return
			}

		case http.MethodPost:
			err := r.ParseMultipartForm(10 << 20)
			if err != nil {
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

			tempFile, err := os.CreateTemp("tmp", "upload-*")
			if err != nil {
				http.Error(w, "Failed to create temp file", http.StatusInternalServerError)
				return
			}
			defer os.Remove(tempFile.Name())
			defer tempFile.Close()

			/////move to seperate function
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
			/////move to seperate function

			if err := os.MkdirAll(finalDir, 0755); err != nil {
				http.Error(w, "Failed to create storage directory", http.StatusInternalServerError)
				return
			}

			content := &models.Content{Title: title, FileMedia: models.MediaType(fileMedia), Filename: finalName}
			err = dao.AddContent(db, content)
			if err != nil {
				http.Error(w, "Failed to insert into DB", http.StatusInternalServerError)
				return
			}

			if err := os.Rename(tempFile.Name(), finalPath); err != nil {
				_ = dao.DeleteContentByFilename(db, finalName)
				http.Error(w, "Failed to store file", http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/view", http.StatusSeeOther)
			return

		default:
			http.Error(w, "Unsupported Method Request", http.StatusMethodNotAllowed)
			return
		}
	})
}

// TODO: add filtered viewing
func handleViewContent(db *sql.DB, tmpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		content, err := dao.GetContentFiles(db)
		if err != nil {
			http.Error(w, "Database query failed", http.StatusInternalServerError)
			return
		}

		user := utils.GetSessionUser(db, r)
		err = tmpl.ExecuteTemplate(w, "view.html", struct {
			Paths []string
			User  *models.User
		}{
			Paths: content,
			User:  user,
		})
		if err != nil {
			http.Error(w, "Template error", http.StatusInternalServerError)
			return
		}
	})
}

func handleViewUploads(db *sql.DB, tmpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("id")
		if err != nil {
			http.Error(w, "Failed to read cookie", http.StatusBadRequest)
			return
		}

		content, err := dao.GetUserContentFiles(db, cookie.Value)
		if err != nil {
			http.Error(w, "Database query failed", http.StatusInternalServerError)
			return
		}

		err = tmpl.ExecuteTemplate(w, "uploads.html", struct{ Paths []string }{
			Paths: content,
		})
		if err != nil {
			http.Error(w, "Template error", http.StatusInternalServerError)
			return
		}
	})
}

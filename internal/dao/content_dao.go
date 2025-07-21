package dao

import (
	"database/sql"
	"goserv/internal/models"
	"log"
	"path"

	"github.com/lib/pq"
)

func AddContent(db *sql.DB, content *models.Content) error {
	_, err := db.Exec("INSERT INTO content (title, media_type, file_name) VALUES ($1, $2, $3)", content.Title, content.FileMedia, content.Filename)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				log.Printf("Tag already exists\n")
				return nil
			} else {
				log.Printf("Error saving tag: %v\n", err)
				return err
			}
		}
	}
	return err
}

func DeleteContentByFilename(db *sql.DB, filename string) error {
	_, err := db.Exec("DELETE FROM content WHERE filename = $1", filename)
	return err
}

// TODO: add filtered viewing
func GetContentFiles(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT file_name FROM content")
	if err != nil {
		log.Printf("Error with database query: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	files := []string{}
	for rows.Next() {
		var file string
		rows.Scan(&file)

		filePath := path.Join(file[0:2], file[2:4], file)

		files = append(files, filePath)
	}
	return files, nil
}

func GetUserContentFiles(db *sql.DB, sessionID string) ([]string, error) {
	rows, err := db.Query(
		"SELECT c.file_name FROM sessions s"+
			" Join users u ON s.username = u.username"+
			" JOIN userposts up ON u.id = up.user_id"+
			" JOIN content c ON up.post_id = c.id"+
			" WHERE s.session_id = $1", sessionID)
	if err != nil {
		log.Printf("Error with database query: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	files := []string{}
	for rows.Next() {
		var file string
		rows.Scan(&file)

		filePath := path.Join(file[0:2], file[2:4], file)

		files = append(files, filePath)
	}
	return files, nil
}

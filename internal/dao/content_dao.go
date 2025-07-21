package dao

import (
	"database/sql"
	"goserv/internal/models"
	"log"
	"path"
)

func AddPost(db *sql.DB, content *models.Content, userID int) error {
	var postID int
	err := db.QueryRow("INSERT INTO Posts (title, media_type, file_name)"+
		" VALUES ($1, $2, $3) RETURNING id",
		content.Title, content.FileMedia, content.Filename).Scan(&postID)
	if err != nil {
		log.Printf("Error saving post: %v\n", err)
		return err
	}

	if err := linkUserToPost(db, userID, postID); err != nil {
		DeletePostByFilename(db, content.Filename)
		return err
	}
	return err
}

func DeletePostByFilename(db *sql.DB, filename string) error {
	_, err := db.Exec("DELETE FROM Posts WHERE filename = $1", filename)
	return err
}

func linkUserToPost(db *sql.DB, userID int, postID int) error {
	_, err := db.Exec("INSERT INTO UserPosts (user_id, post_id) VALUES ($1, $2)", userID, postID)
	if err != nil {
		log.Printf("Error linking user to post: %v\n", err)
		return err
	}
	return err
}

// TODO: add filtered viewing
func GetContentFiles(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT filename FROM Posts")
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
		"SELECT p.filename FROM Sessions s"+
			" Join Users u ON s.username = u.username"+
			" JOIN UserPosts up ON u.id = up.user_id"+
			" JOIN Posts p ON up.post_id = p.id"+
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

func GetUserFavContentFiles(db *sql.DB, sessionID string) ([]string, error) {
	rows, err := db.Query(
		"SELECT p.filename FROM Sessions s"+
			" Join Users u ON s.username = u.username"+
			" JOIN UserFavs uf ON u.id = uf.user_id"+
			" JOIN Posts p ON uf.post_id = p.id"+
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

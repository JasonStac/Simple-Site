package dao

import (
	"database/sql"
	"goserv/internal/models"
	"log"
	"path"

	"github.com/lib/pq"
)

type ContentDao struct {
	db *sql.DB
}

func NewContentDao(db *sql.DB) *ContentDao {
	return &ContentDao{db: db}
}

func (dao *ContentDao) AddContent(content *models.Content) error {
	_, err := dao.db.Exec("INSERT INTO content (title, media_type, file_name, artist) VALUES ($1, $2, $3, $4)", content.Title, content.FileMedia, content.Filename, content.Artist)
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

func (dao *ContentDao) DeleteContentByFilename(filename string) error {
	_, err := dao.db.Exec("DELETE FROM content WHERE filename = $1", filename)
	return err
}

// TODO: add filtered viewing
func (dao *ContentDao) GetContentFiles() ([]string, error) {
	rows, err := dao.db.Query("SELECT file_name FROM content")
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

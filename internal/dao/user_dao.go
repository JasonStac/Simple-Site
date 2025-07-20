package dao

import (
	"database/sql"
	"goserv/internal/models"
	"log"
)

func AddUser(db *sql.DB, user *models.User) error {
	_, err := db.Exec("INSERT INTO users (username, pass_hash) VALUES ($1, $2)", user.Username, user.PassHash)
	if err != nil {
		log.Printf("Error adding user: %v\n", err)
	}
	return err
}

func GetUser(db *sql.DB, username string) (*models.User, error) {
	var user = &models.User{}
	err := db.QueryRow("SELECT * FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.PassHash, &user.IsAdmin)
	if err != nil {
		log.Printf("Error getting user: %v\n", err)
		return nil, err
	}
	return user, nil
}

func SaveSession(db *sql.DB, username string, sessionID string) error {
	_, err := db.Exec("INSERT INTO sessions (username, session_id) VALUES ($1, $2) ON CONFLICT (username) DO UPDATE SET session_id = $2", username, sessionID)
	if err != nil {
		log.Printf("Error saving session: %v\n", err)
	}
	return err
}

func GetSessionUsername(db *sql.DB, sessionID string) (string, error) {
	var username string
	err := db.QueryRow("SELECT username FROM sessions WHERE session_id = $1", sessionID).Scan(&username)
	if err != nil {
		log.Printf("Error getting uesr: %v\n", err)
		return "", err
	}
	return username, nil
}

func DeleteSession(db *sql.DB, username string) error {
	_, err := db.Exec("DELETE FROM sessions WHERE username = $1", username)
	return err
}

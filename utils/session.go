package utils

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"goserv/internal/dao"
	"goserv/internal/models"
	"net/http"
)

func GenerateSessionID() string {
	b := make([]byte, 64)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func GetSessionUser(db *sql.DB, r *http.Request) *models.User {
	var user *models.User
	cookie, err := r.Cookie("id")
	if err != nil {
		return user
	}

	username, err := dao.GetSessionUsername(db, cookie.Value)
	if err != nil {
		return user
	}

	user, err = dao.GetUser(db, username)
	if err != nil {
		return user
	}
	return user
}

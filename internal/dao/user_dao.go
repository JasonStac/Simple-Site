package dao

import (
	"database/sql"
	"goserv/internal/models"
	"log"
)

type UserDao struct {
	db *sql.DB
}

func NewUserDao(db *sql.DB) *UserDao {
	return &UserDao{db: db}
}

func (dao *UserDao) AddUser(user *models.User) error {
	_, err := dao.db.Exec("INSERT INTO users (username, pass_hash) VALUES ($1, $2)", user.Username, user.PassHash)
	if err != nil {
		log.Printf("Error adding user: %v\n", err)
	}

	return err
}

func (dao *UserDao) GetUser(username string) (*models.User, error) {
	var user = &models.User{}
	err := dao.db.QueryRow("SELECT * FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.PassHash, &user.IsAdmin)
	if err != nil {
		log.Printf("Error with getting user: %v\n", err)
		return nil, err
	}

	return user, nil
}

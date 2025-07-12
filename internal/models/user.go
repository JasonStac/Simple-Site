package models

type User struct {
	ID       int
	Username string
	PassHash string
	IsAdmin  bool
}

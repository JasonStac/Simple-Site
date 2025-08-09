package service

import (
	"context"
	"goserv/internal/domain/users"
	"goserv/internal/domain/users/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetByUsername(ctx context.Context, username string) (*users.User, error) {
	return s.repo.GetByUsername(ctx, username)
}

func (s *UserService) CheckPassword(ctx context.Context, username string, password string) (*users.User, bool, error) {
	return s.repo.CheckPassword(ctx, username, password)
}

func (s *UserService) Register(ctx context.Context, username string, password string) error {
	// TODO: test if user exists before trying to make new one
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &users.User{
		Username: username,
		IsAdmin:  false,
	}

	return s.repo.Register(ctx, user, string(hashedPass))
}

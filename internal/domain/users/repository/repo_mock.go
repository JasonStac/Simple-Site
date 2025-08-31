package repository

import (
	"context"
	"goserv/internal/domain/users"
)

type UserMock struct {
	RegisterFunc      func(ctx context.Context, user *users.User, passHash string) error
	GetByUsernameFunc func(ctx context.Context, username string) (*users.User, error)
	CheckPasswordFunc func(ctx context.Context, username string, password string) (*users.User, bool, error)
	GetByUserIDFunc   func(ctx context.Context, userID int) (*users.User, error)
	IsAdminFunc       func(ctx context.Context, userID int) (bool, error)
}

func (m *UserMock) Register(ctx context.Context, user *users.User, passHash string) error {
	return m.RegisterFunc(ctx, user, passHash)
}

func (m *UserMock) GetByUsername(ctx context.Context, username string) (*users.User, error) {
	return m.GetByUsernameFunc(ctx, username)
}

func (m *UserMock) CheckPassword(ctx context.Context, username string, password string) (*users.User, bool, error) {
	return m.CheckPasswordFunc(ctx, username, password)
}

func (m *UserMock) GetByUserID(ctx context.Context, userID int) (*users.User, error) {
	return m.GetByUserIDFunc(ctx, userID)
}

func (m *UserMock) IsAdmin(ctx context.Context, userID int) (bool, error) {
	return m.IsAdminFunc(ctx, userID)
}

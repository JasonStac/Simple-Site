package repository

import (
	"context"
	"goserv/ent/gen"
	entUser "goserv/ent/gen/user"
	"goserv/internal/domain/users"

	"golang.org/x/crypto/bcrypt"
)

type User interface {
	Register(ctx context.Context, user *users.User, passHash string) error
	GetByUsername(ctx context.Context, username string) (*users.User, error)
	CheckPassword(ctx context.Context, username string, password string) (*users.User, bool, error)
	GetByUserID(ctx context.Context, userID int) (*users.User, error)
	IsAdmin(ctx context.Context, userID int) (bool, error)
}

type userRepository struct {
	client *gen.Client
}

func NewUserRepository(client *gen.Client) *userRepository {
	return &userRepository{client: client}
}

func (repo *userRepository) Register(ctx context.Context, user *users.User, passHash string) error {
	_, err := repo.client.User.Create().SetUsername(user.Username).SetPassHash(passHash).SetIsAdmin(user.IsAdmin).Save(ctx)
	return err
}

func (repo *userRepository) GetByUsername(ctx context.Context, username string) (*users.User, error) {
	user, err := repo.client.User.Query().Where(entUser.UsernameEQ(username)).Only(ctx)
	if err != nil {
		return nil, err
	}
	return &users.User{ID: user.ID, Username: user.Username, IsAdmin: user.IsAdmin}, nil
}

func (repo *userRepository) CheckPassword(ctx context.Context, username string, password string) (*users.User, bool, error) {
	user, err := repo.client.User.Query().Where(entUser.UsernameEQ(username)).Only(ctx)
	if err != nil {
		return nil, false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(password))
	if err != nil {
		return nil, false, err
	}
	return &users.User{ID: user.ID, Username: user.Username, IsAdmin: user.IsAdmin}, true, nil
}

func (repo *userRepository) GetByUserID(ctx context.Context, userID int) (*users.User, error) {
	user, err := repo.client.User.Query().Where(entUser.IDEQ(userID)).Only(ctx)
	if err != nil {
		return nil, err
	}
	return &users.User{ID: user.ID, Username: user.Username, IsAdmin: user.IsAdmin}, nil
}

func (repo *userRepository) IsAdmin(ctx context.Context, userID int) (bool, error) {
	user, err := repo.client.User.Query().Select(entUser.FieldIsAdmin).Where(entUser.IDEQ(userID)).Only(ctx)
	if err != nil {
		return false, err
	}
	return user.IsAdmin, nil
}

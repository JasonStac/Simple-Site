package repository

import (
	"context"
	"goserv/ent/gen"
	"goserv/ent/gen/user"
	"goserv/internal/domain/users"

	"golang.org/x/crypto/bcrypt"
)

type User interface {
	Register(ctx context.Context, user *users.User, passHash string) error
	GetByUsername(ctx context.Context, username string) (*users.User, error)
	CheckPassword(ctx context.Context, username string, password string) (*users.User, bool, error)
	GetByUserID(ctx context.Context, userID int) (*users.User, error)
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
	user, err := repo.client.User.Query().Where(user.UsernameEQ(username)).Only(ctx)
	return &users.User{ID: user.ID, Username: user.Username, IsAdmin: user.IsAdmin}, err
}

func (repo *userRepository) CheckPassword(ctx context.Context, username string, password string) (*users.User, bool, error) {
	user, err := repo.client.User.Query().Where(user.UsernameEQ(username)).Only(ctx)
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
	user, err := repo.client.User.Query().Where(user.IDEQ(userID)).Only(ctx)
	return &users.User{ID: user.ID, Username: user.Username, IsAdmin: user.IsAdmin}, err
}

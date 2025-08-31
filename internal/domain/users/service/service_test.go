package service

import (
	"context"
	"errors"
	"goserv/internal/domain/users"
	"goserv/internal/domain/users/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserService_GetByUsername(t *testing.T) {
	type args struct {
		name string
	}
	type want struct {
		user *users.User
		err  error
	}
	type test struct {
		name string
		args args
		want want
	}

	tests := []test{
		{
			name: "simple get",
			args: args{
				name: "username",
			},
			want: want{
				user: &users.User{
					ID:       1,
					Username: "username",
					IsAdmin:  false,
				},
				err: nil,
			},
		},
		{
			name: "no user",
			args: args{
				name: "username",
			},
			want: want{
				user: nil,
				err:  errors.New("test error"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			userRepo := &repository.UserMock{
				GetByUsernameFunc: func(ctx context.Context, username string) (*users.User, error) {
					return test.want.user, test.want.err
				},
			}

			service := NewUserService(userRepo)

			user, err := service.GetByUsername(context.Background(), test.args.name)
			assert.Equal(t, test.want.err, err)
			assert.Equal(t, test.want.user, user)
		})
	}
}

func TestUserService_CheckPassword(t *testing.T) {
	type args struct {
		name     string
		password string
	}
	type want struct {
		user    *users.User
		isMatch bool
		err     error
	}
	type test struct {
		name string
		args args
		want want
	}

	basicArgs := args{
		name:     "username",
		password: "password",
	}

	tests := []test{
		{
			name: "simple match",
			args: basicArgs,
			want: want{
				user: &users.User{
					ID:       1,
					Username: "username",
					IsAdmin:  false,
				},
				isMatch: true,
				err:     nil,
			},
		},
		{
			name: "no match",
			args: basicArgs,
			want: want{
				user:    nil,
				isMatch: false,
				err:     nil,
			},
		},
		{
			name: "error on match",
			args: basicArgs,
			want: want{
				user:    nil,
				isMatch: false,
				err:     errors.New("test error"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			userRepo := &repository.UserMock{
				CheckPasswordFunc: func(ctx context.Context, username string, password string) (*users.User, bool, error) {
					return test.want.user, test.want.isMatch, test.want.err
				},
			}

			service := NewUserService(userRepo)

			user, isMatch, err := service.CheckPassword(context.Background(), test.args.name, test.args.password)
			assert.Equal(t, test.want.err, err)
			assert.Equal(t, test.want.isMatch, isMatch)
			assert.Equal(t, test.want.user, user)
		})
	}
}

func TestUserService_Register(t *testing.T) {
	type args struct {
		username string
		password string
	}
	type want struct {
		err error
	}
	type test struct {
		name string
		args args
		want want
	}

	basicArgs := args{
		username: "username",
		password: "password",
	}

	tests := []test{
		{
			name: "simple register",
			args: basicArgs,
			want: want{
				err: nil,
			},
		},
		{
			name: "error registering",
			args: basicArgs,
			want: want{
				err: errors.New("test error"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			userRepo := &repository.UserMock{
				RegisterFunc: func(ctx context.Context, user *users.User, passHash string) error {
					return test.want.err
				},
			}

			service := NewUserService(userRepo)

			err := service.Register(context.Background(), test.args.username, test.args.password)
			assert.Equal(t, test.want.err, err)
		})
	}
}

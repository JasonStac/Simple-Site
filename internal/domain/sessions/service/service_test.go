package service

import (
	"context"
	"errors"
	"goserv/internal/domain/sessions"
	"goserv/internal/domain/sessions/repository"
	"goserv/internal/domain/users"
	uRepo "goserv/internal/domain/users/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSessionService_Login(t *testing.T) {
	type args struct {
		username string
		password string
	}
	type want struct {
		user     *users.User
		isMatch  bool
		checkErr error
		err      error
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

	basicUser := users.User{
		ID:       1,
		Username: "username",
		IsAdmin:  false,
	}

	tests := []test{
		{
			name: "simple login",
			args: basicArgs,
			want: want{
				user:     &basicUser,
				isMatch:  true,
				checkErr: nil,
				err:      nil,
			},
		},
		{
			name: "bad password",
			args: basicArgs,
			want: want{
				user:     nil,
				isMatch:  false,
				checkErr: errors.New("invalid credentials"),
				err:      nil,
			},
		},
		{
			name: "login error",
			args: basicArgs,
			want: want{
				user:     &basicUser,
				isMatch:  true,
				checkErr: nil,
				err:      errors.New("test error"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sessionRepo := &repository.SessionMock{
				LoginFunc: func(ctx context.Context, user *sessions.Session) error {
					return test.want.err
				},
			}

			userRepo := &uRepo.UserMock{
				CheckPasswordFunc: func(ctx context.Context, username string, password string) (*users.User, bool, error) {
					return test.want.user, test.want.isMatch, test.want.checkErr
				},
			}

			service := NewSessionService(sessionRepo, userRepo)

			sessionID, err := service.Login(context.Background(), test.args.username, test.args.password)
			if test.want.checkErr == nil {
				assert.Equal(t, test.want.err, err)
			} else {
				assert.Equal(t, test.want.checkErr, err)
			}
			if err == nil {
				assert.NotEmpty(t, sessionID)
			}
		})
	}
}

func TestSessionService_Logout(t *testing.T) {
	type args struct {
		sessionID string
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
		sessionID: "session",
	}

	tests := []test{
		{
			name: "simple logout",
			args: basicArgs,
			want: want{
				err: nil,
			},
		},
		{
			name: "error logout",
			args: basicArgs,
			want: want{
				err: errors.New("test error"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sessionRepo := &repository.SessionMock{
				LogoutFunc: func(ctx context.Context, sessionID string) error {
					return test.want.err
				},
			}

			service := NewSessionService(sessionRepo, nil)

			err := service.Logout(context.Background(), test.args.sessionID)
			assert.Equal(t, test.want.err, err)
		})
	}
}

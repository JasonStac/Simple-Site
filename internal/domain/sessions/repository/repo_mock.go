package repository

import (
	"context"
	"goserv/internal/domain/sessions"
)

type SessionMock struct {
	LoginFunc                func(ctx context.Context, session *sessions.Session) error
	LogoutFunc               func(ctx context.Context, sessionID string) error
	GetUserIDBySessionIDFunc func(ctx context.Context, sessionID string) (int, error)
}

func (m *SessionMock) Login(ctx context.Context, session *sessions.Session) error {
	return m.LoginFunc(ctx, session)
}

func (m *SessionMock) Logout(ctx context.Context, sessionID string) error {
	return m.LogoutFunc(ctx, sessionID)
}

func (m *SessionMock) GetUserIDBySessionID(ctx context.Context, sessionID string) (int, error) {
	return m.GetUserIDBySessionIDFunc(ctx, sessionID)
}

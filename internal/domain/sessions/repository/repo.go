package repository

import (
	"context"
	"goserv/ent/gen"
	entSession "goserv/ent/gen/session"
	"goserv/internal/domain/sessions"
)

type Session interface {
	Login(ctx context.Context, session *sessions.Session) error
	Logout(ctx context.Context, sessionID string) error
	//GetUsernameBySessionID(ctx context.Context, sessionID string) (string, error)
	GetUserIDBySessionID(ctx context.Context, sessionID string) (int, error)
}

type sessionRepo struct {
	client *gen.Client
}

func NewSessionRepository(client *gen.Client) *sessionRepo {
	return &sessionRepo{client: client}
}

func (repo *sessionRepo) Login(ctx context.Context, session *sessions.Session) error {
	_, err := repo.client.Session.Create().SetID(session.ID).SetUserID(session.UserID).Save(ctx)
	return err
}

func (repo *sessionRepo) Logout(ctx context.Context, sessionID string) error {
	_, err := repo.client.Session.Delete().Where(entSession.IDEQ(sessionID)).Exec(ctx)
	return err
}

// func (repo *sessionRepo) GetUsernameBySessionID(ctx context.Context, sessionID string) (string, error) {
// 	var username string
// 	err := repo.db.Where("session_id = ?", sessionID).First(&username).Error
// 	return username, err
// }

func (repo *sessionRepo) GetUserIDBySessionID(ctx context.Context, sessionID string) (int, error) {
	session, err := repo.client.Session.Query().Where(entSession.IDEQ(sessionID)).Only(ctx)
	if err != nil {
		return 0, err
	}
	return session.UserID, nil
}

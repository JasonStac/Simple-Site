package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"goserv/internal/domain/sessions"
	"goserv/internal/domain/sessions/repository"
	"goserv/internal/domain/users/service"
)

type SessionService struct {
	repo    repository.Session
	userSvc *service.UserService
}

func NewSessionService(repo repository.Session, userSvc *service.UserService) *SessionService {
	return &SessionService{repo: repo, userSvc: userSvc}
}

func (s *SessionService) Login(ctx context.Context, username string, password string) (string, error) {
	user, ok, err := s.userSvc.CheckPassword(ctx, username, password)
	if err != nil || !ok {
		return "", errors.New("invalid credentials")
	}

	// generateID and save session
	sessionID := generateSessionID()
	session := &sessions.Session{ID: sessionID, UserID: user.ID}
	err = s.repo.Login(ctx, session)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func (s *SessionService) Logout(ctx context.Context, sessionID string) error {
	return s.repo.Logout(ctx, sessionID)
}

func generateSessionID() string {
	b := make([]byte, 64)
	rand.Read(b)
	return hex.EncodeToString(b)
}

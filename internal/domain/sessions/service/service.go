package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"goserv/internal/domain/sessions"
	"goserv/internal/domain/sessions/repository"
	uRepo "goserv/internal/domain/users/repository"
)

type SessionService struct {
	repo     repository.Session
	userRepo uRepo.User
}

func NewSessionService(repo repository.Session, userRepo uRepo.User) *SessionService {
	return &SessionService{repo: repo, userRepo: userRepo}
}

func (s *SessionService) Login(ctx context.Context, username string, password string) (string, error) {
	user, isMatch, err := s.userRepo.CheckPassword(ctx, username, password)
	if err != nil || !isMatch {
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

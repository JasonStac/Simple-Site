package middleware

import (
	"context"
	"goserv/internal/domain/sessions/repository"
	"net/http"
)

type key string

const userKey key = "user_id"

func AuthRestrictMiddleware(sessionRepo repository.Session) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("id")
			if err != nil {
				http.Error(w, "Unauthorized: no session", http.StatusUnauthorized)
				return
			}

			userID, err := sessionRepo.GetUserIDBySessionID(r.Context(), cookie.Value)
			if err != nil {
				http.Error(w, "Unauthorized: invalid session", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AuthCheckMiddleware(sessionRepo repository.Session) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := -1
			cookie, err := r.Cookie("id")
			if err == nil {
				userID, err = sessionRepo.GetUserIDBySessionID(r.Context(), cookie.Value)
				if err != nil {
					userID = -1
				}
			}

			ctx := context.WithValue(r.Context(), userKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(r *http.Request) (int, bool) {
	username, ok := r.Context().Value(userKey).(int)
	return username, ok
}

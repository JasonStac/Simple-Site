package middleware

import (
	"context"
	pRepo "goserv/internal/domain/posts/repository"
	uRepo "goserv/internal/domain/users/repository"
	"net/http"
	"strconv"
)

func DeleteMiddleware(userRepo uRepo.User, postRepo pRepo.Post) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := GetUserID(r)
			if !ok {
				http.Error(w, "Error getting user id", http.StatusInternalServerError)
				return
			}

			postID, err := strconv.Atoi(r.FormValue("id"))
			if err != nil {
				http.Error(w, "Error getting post id", http.StatusBadRequest)
				return
			}
			ctx := context.WithValue(r.Context(), postKey, postID)

			post, err := postRepo.GetPost(r.Context(), postID)
			if err != nil {
				http.Error(w, "Error deleting", http.StatusInternalServerError)
				return
			}
			ctx = context.WithValue(ctx, filenameKey, post.Filename)
			ctx = context.WithValue(ctx, fileExtKey, post.FileExt)

			isAdmin, err := userRepo.IsAdmin(r.Context(), userID)
			if err != nil {
				http.Error(w, "Error deleting", http.StatusInternalServerError)
				return
			}

			if isAdmin {
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			if post.OwnerID == userID {
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		})
	}
}

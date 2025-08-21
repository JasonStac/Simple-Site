package middleware

import "net/http"

type key string

const userKey key = "user_id"
const filepathKey key = "filepath"
const postKey key = "post_id"

func GetUserID(r *http.Request) (int, bool) {
	userID, ok := r.Context().Value(userKey).(int)
	return userID, ok
}

func GetFilepath(r *http.Request) (string, bool) {
	filepath, ok := r.Context().Value(filepathKey).(string)
	return filepath, ok
}

func GetPostID(r *http.Request) (int, bool) {
	postID, ok := r.Context().Value(postKey).(int)
	return postID, ok
}

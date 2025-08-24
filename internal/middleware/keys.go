package middleware

import (
	"goserv/internal/domain/tags"
	"net/http"
)

type key string

const userKey key = "user_id"
const filenameKey key = "filename"
const fileExtKey key = "file_ext"
const postKey key = "post_id"
const tagKey key = "tags"

func GetUserID(r *http.Request) (int, bool) {
	userID, ok := r.Context().Value(userKey).(int)
	return userID, ok
}

func GetCompleteFilename(r *http.Request) (string, string, bool) {
	filename, ok1 := r.Context().Value(filenameKey).(string)
	fileExt, ok2 := r.Context().Value(fileExtKey).(string)
	return filename, fileExt, ok1 && ok2
}

func GetPostID(r *http.Request) (int, bool) {
	postID, ok := r.Context().Value(postKey).(int)
	return postID, ok
}

func GetTags(r *http.Request) ([]tags.Tag, bool) {
	tags, ok := r.Context().Value(tagKey).([]tags.Tag)
	return tags, ok
}

package utils

import (
	"strconv"
	"strings"
)

func GetPostIDFromPath(url string) (int, error) {
	return strconv.Atoi(strings.TrimPrefix(url, "/view/posts/"))
}

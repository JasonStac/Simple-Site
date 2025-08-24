package posts

import (
	"goserv/internal/domain/tags"
	"goserv/internal/static/enum"
)

type Post struct {
	ID        int
	Title     string
	MediaType enum.MediaType
	Filename  string
	FileExt   string
	OwnerID   int

	Tags []tags.Tag
}

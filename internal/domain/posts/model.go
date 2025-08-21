package posts

import (
	"goserv/internal/domain/artists"
	"goserv/internal/domain/tags"
	"goserv/internal/models"
)

type Post struct {
	ID        int
	Title     string
	MediaType models.MediaType
	Filename  string
	OwnerID   int

	Artists []artists.Artist
	Tags    []tags.Tag
}

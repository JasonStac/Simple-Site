package tags

import (
	"goserv/internal/static/enum"
)

type Tag struct {
	ID   int          `json:"id,omitempty"`
	Type enum.TagType `json:"type"`
	Name string       `json:"value"`
}

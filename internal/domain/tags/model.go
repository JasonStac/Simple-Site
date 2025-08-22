package tags

import (
	"goserv/internal/models"
)

type Tag struct {
	ID   int            `json:"id,omitempty"`
	Type models.TagType `json:"type"`
	Name string         `json:"value"`
}

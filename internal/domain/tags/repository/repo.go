package repository

import (
	"context"
	"goserv/ent/gen"
	"goserv/internal/domain/tags"
)

type Tag interface {
	ListTags(ctx context.Context) ([]tags.Tag, error)
}

type tagRepository struct {
	client *gen.Client
}

func NewTagRepository(client *gen.Client) *tagRepository {
	return &tagRepository{client: client}
}

func (repo *tagRepository) ListTags(ctx context.Context) ([]tags.Tag, error) {
	entTags, err := repo.client.Tag.Query().All(ctx)
	returnTags := make([]tags.Tag, len(entTags))
	for i, t := range entTags {
		returnTags[i] = tags.Tag{ID: t.ID, Name: t.Name}
	}
	return returnTags, err
}

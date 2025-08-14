package repository

import (
	"context"
	"goserv/ent/gen"
	"goserv/internal/domain/tags"
)

type Tag interface {
	AddTag(ctx context.Context, name string) (int, error)
	ListTags(ctx context.Context) ([]tags.Tag, error)
}

type tagRepository struct {
	client *gen.Client
}

func NewTagRepository(client *gen.Client) *tagRepository {
	return &tagRepository{client: client}
}

func (repo *tagRepository) AddTag(ctx context.Context, name string) (int, error) {
	entTag, err := repo.client.Tag.Create().Save(ctx)
	if err != nil {
		return -1, err
	}
	return entTag.ID, nil
}

func (repo *tagRepository) ListTags(ctx context.Context) ([]tags.Tag, error) {
	entTags, err := repo.client.Tag.Query().All(ctx)
	returnTags := make([]tags.Tag, len(entTags))
	for i, t := range entTags {
		returnTags[i] = tags.Tag{ID: t.ID, Name: t.Name}
	}
	return returnTags, err
}

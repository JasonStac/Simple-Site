package repository

import (
	"context"
	"goserv/ent/gen"
	entTag "goserv/ent/gen/tag"
	"goserv/internal/domain/tags"
	"goserv/internal/models"
)

type Tag interface {
	AddTag(ctx context.Context, name string, tagType models.TagType) (int, error)
	ListTags(ctx context.Context) ([]tags.Tag, error)
	ListGeneralTags(ctx context.Context) ([]tags.Tag, error)
	ListPeopleTags(ctx context.Context) ([]tags.Tag, error)
}

type tagRepository struct {
	client *gen.Client
}

func NewTagRepository(client *gen.Client) *tagRepository {
	return &tagRepository{client: client}
}

func (repo *tagRepository) AddTag(ctx context.Context, name string, tagType models.TagType) (int, error) {
	entTag, err := repo.client.Tag.Create().SetName(name).SetTagType(entTag.TagType(tagType)).Save(ctx)
	if err != nil {
		return 0, err
	}
	return entTag.ID, nil
}

func (repo *tagRepository) ListTags(ctx context.Context) ([]tags.Tag, error) {
	entTags, err := repo.client.Tag.Query().All(ctx)
	returnTags := make([]tags.Tag, len(entTags))
	for i := range entTags {
		returnTags[i] = tags.Tag{
			ID:   entTags[i].ID,
			Type: models.TagType(entTags[i].TagType),
			Name: entTags[i].Name,
		}
	}
	return returnTags, err
}

func (repo *tagRepository) ListGeneralTags(ctx context.Context) ([]tags.Tag, error) {
	entTags, err := repo.client.Tag.Query().Where(entTag.TagTypeEQ(entTag.TagTypeGeneral)).All(ctx)
	if err != nil {
		return nil, err
	}
	returnTags := make([]tags.Tag, len(entTags))
	for i := range entTags {
		returnTags[i] = tags.Tag{
			ID:   entTags[i].ID,
			Type: models.TagGeneral,
			Name: entTags[i].Name,
		}
	}
	return returnTags, err
}

func (repo *tagRepository) ListPeopleTags(ctx context.Context) ([]tags.Tag, error) {
	entTags, err := repo.client.Tag.Query().Where(entTag.TagTypeEQ(entTag.TagTypePeople)).All(ctx)
	if err != nil {
		return nil, err
	}
	returnTags := make([]tags.Tag, len(entTags))
	for i := range entTags {
		returnTags[i] = tags.Tag{
			ID:   entTags[i].ID,
			Type: models.TagPeople,
			Name: entTags[i].Name,
		}
	}
	return returnTags, err
}

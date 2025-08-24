package service

import (
	"context"
	"errors"
	"goserv/internal/domain/tags"
	"goserv/internal/domain/tags/repository"
	"goserv/internal/static/enum"
)

type TagService struct {
	repo repository.Tag
}

func NewTagService(repo repository.Tag) *TagService {
	return &TagService{repo: repo}
}

func (s *TagService) AddTag(ctx context.Context, name string, tagType enum.TagType) (int, error) {
	return s.repo.AddTag(ctx, name, tagType)
}

func (s *TagService) ListTags(ctx context.Context) ([]tags.Tag, error) {
	return s.repo.ListTags(ctx)
}

func (s *TagService) ListGeneralTags(ctx context.Context) ([]tags.Tag, error) {
	return s.repo.ListGeneralTags(ctx)
}

func (s *TagService) ListPeopleTags(ctx context.Context) ([]tags.Tag, error) {
	return s.repo.ListPeopleTags(ctx)
}

func (s *TagService) SeperateTagTypes(ctx context.Context, allTags []tags.Tag) (map[enum.TagType][]tags.Tag, error) {
	if allTags == nil {
		return nil, nil
	}

	result := make(map[enum.TagType][]tags.Tag, len(allTags))

	for i := range allTags {
		switch enum.TagType(allTags[i].Type) {
		case enum.TagGeneral, enum.TagPeople:
			result[enum.TagType(allTags[i].Type)] = append(result[enum.TagType(allTags[i].Type)], allTags[i])
		default:
			return nil, errors.New("invalid tag type detected")
		}
	}
	return result, nil
}

package service

import (
	"context"
	"goserv/internal/domain/tags"
	"goserv/internal/domain/tags/repository"
)

type TagService struct {
	repo repository.Tag
}

func NewTagService(repo repository.Tag) *TagService {
	return &TagService{repo: repo}
}

func (s *TagService) ListTags(ctx context.Context) ([]tags.Tag, error) {
	return s.repo.ListTags(ctx)
}

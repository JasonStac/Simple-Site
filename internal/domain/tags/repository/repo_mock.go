package repository

import (
	"context"
	"goserv/internal/domain/tags"
	"goserv/internal/static/enum"
)

type TagMock struct {
	AddTagFunc          func(ctx context.Context, name string, tagType enum.TagType) (int, error)
	ListTagsFunc        func(ctx context.Context) ([]tags.Tag, error)
	ListGeneralTagsFunc func(ctx context.Context) ([]tags.Tag, error)
	ListPeopleTagsFunc  func(ctx context.Context) ([]tags.Tag, error)
}

func (m *TagMock) AddTag(ctx context.Context, name string, tagType enum.TagType) (int, error) {
	return m.AddTagFunc(ctx, name, tagType)
}

func (m *TagMock) ListTags(ctx context.Context) ([]tags.Tag, error) {
	return m.ListTagsFunc(ctx)
}

func (m *TagMock) ListGeneralTags(ctx context.Context) ([]tags.Tag, error) {
	return m.ListGeneralTagsFunc(ctx)
}

func (m *TagMock) ListPeopleTags(ctx context.Context) ([]tags.Tag, error) {
	return m.ListPeopleTagsFunc(ctx)
}

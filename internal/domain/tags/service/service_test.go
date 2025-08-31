package service

import (
	"context"
	"errors"
	"goserv/internal/domain/tags"
	"goserv/internal/domain/tags/repository"
	"goserv/internal/static/enum"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTagService_AddTag(t *testing.T) {
	type args struct {
		name    string
		tagType enum.TagType
	}
	type want struct {
		tagID int
		err   error
	}
	type test struct {
		name string
		args args
		want want
	}

	tests := []test{
		{
			name: "simple general add",
			args: args{
				name:    "general tag",
				tagType: enum.TagGeneral,
			},
			want: want{
				tagID: 1,
				err:   nil,
			},
		},
		{
			name: "simple people add",
			args: args{
				name:    "people tag",
				tagType: enum.TagPeople,
			},
			want: want{
				tagID: 1,
				err:   nil,
			},
		},
		{
			name: "failed tag add",
			args: args{
				name:    "tag",
				tagType: enum.TagGeneral,
			},
			want: want{
				tagID: 0,
				err:   errors.New("test error"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tagRepo := &repository.TagMock{
				AddTagFunc: func(ctx context.Context, name string, tagType enum.TagType) (int, error) {
					return test.want.tagID, test.want.err
				},
			}

			service := NewTagService(tagRepo)

			tagID, err := service.AddTag(context.Background(), test.args.name, test.args.tagType)
			assert.Equal(t, test.want.err, err)
			assert.Equal(t, test.want.tagID, tagID)
		})
	}
}

func TestTagService_ListTags(t *testing.T) {
	type want struct {
		tags []tags.Tag
		err  error
	}
	type test struct {
		name string
		want want
	}

	multipleTags := []tags.Tag{
		{
			ID:   1,
			Type: enum.TagGeneral,
			Name: "general tag",
		},
		{
			ID:   2,
			Type: enum.TagPeople,
			Name: "people tag",
		},
	}

	singleTag := []tags.Tag{
		{
			ID:   1,
			Type: enum.TagGeneral,
			Name: "general tag",
		},
	}

	emptyTags := []tags.Tag{}

	tests := []test{
		{
			name: "multiple tag list",
			want: want{
				tags: multipleTags,
				err:  nil,
			},
		},
		{
			name: "single tag list",
			want: want{
				tags: singleTag,
				err:  nil,
			},
		},
		{
			name: "empty tag list",
			want: want{
				tags: emptyTags,
				err:  errors.New("test error"),
			},
		},
		{
			name: "failed tag list",
			want: want{
				tags: nil,
				err:  errors.New("test error"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tagRepo := &repository.TagMock{
				ListTagsFunc: func(ctx context.Context) ([]tags.Tag, error) {
					return test.want.tags, test.want.err
				},
			}

			service := NewTagService(tagRepo)

			tags, err := service.ListTags(context.Background())
			assert.Equal(t, test.want.err, err)
			assert.Equal(t, test.want.tags, tags)
		})
	}
}

func TestTagService_ListGeneralTags(t *testing.T) {
	type want struct {
		tags []tags.Tag
		err  error
	}
	type test struct {
		name string
		want want
	}

	multipleTags := []tags.Tag{
		{
			ID:   1,
			Type: enum.TagGeneral,
			Name: "general tag",
		},
		{
			ID:   2,
			Type: enum.TagGeneral,
			Name: "general tag2",
		},
	}

	singleTag := []tags.Tag{
		{
			ID:   1,
			Type: enum.TagGeneral,
			Name: "general tag",
		},
	}

	emptyTags := []tags.Tag{}

	tests := []test{
		{
			name: "multiple general tag list",
			want: want{
				tags: multipleTags,
				err:  nil,
			},
		},
		{
			name: "single general tag list",
			want: want{
				tags: singleTag,
				err:  nil,
			},
		},
		{
			name: "empty general tag list",
			want: want{
				tags: emptyTags,
				err:  errors.New("test error"),
			},
		},
		{
			name: "failed general tag list",
			want: want{
				tags: nil,
				err:  errors.New("test error"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tagRepo := &repository.TagMock{
				ListGeneralTagsFunc: func(ctx context.Context) ([]tags.Tag, error) {
					return test.want.tags, test.want.err
				},
			}

			service := NewTagService(tagRepo)

			tags, err := service.ListGeneralTags(context.Background())
			assert.Equal(t, test.want.err, err)
			assert.Equal(t, test.want.tags, tags)
		})
	}
}

func TestTagService_ListPeopleTags(t *testing.T) {
	type want struct {
		tags []tags.Tag
		err  error
	}
	type test struct {
		name string
		want want
	}

	multipleTags := []tags.Tag{
		{
			ID:   1,
			Type: enum.TagPeople,
			Name: "people tag",
		},
		{
			ID:   2,
			Type: enum.TagPeople,
			Name: "people tag2",
		},
	}

	singleTag := []tags.Tag{
		{
			ID:   1,
			Type: enum.TagPeople,
			Name: "people tag",
		},
	}

	emptyTags := []tags.Tag{}

	tests := []test{
		{
			name: "multiple people tag list",
			want: want{
				tags: multipleTags,
				err:  nil,
			},
		},
		{
			name: "single people tag list",
			want: want{
				tags: singleTag,
				err:  nil,
			},
		},
		{
			name: "empty people tag list",
			want: want{
				tags: emptyTags,
				err:  errors.New("test error"),
			},
		},
		{
			name: "failed people tag list",
			want: want{
				tags: nil,
				err:  errors.New("test error"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tagRepo := &repository.TagMock{
				ListPeopleTagsFunc: func(ctx context.Context) ([]tags.Tag, error) {
					return test.want.tags, test.want.err
				},
			}

			service := NewTagService(tagRepo)

			tags, err := service.ListPeopleTags(context.Background())
			assert.Equal(t, test.want.err, err)
			assert.Equal(t, test.want.tags, tags)
		})
	}
}

func TestTagService_SeperateTagTypes(t *testing.T) {
	type args struct {
		allTags []tags.Tag
	}
	type want struct {
		splitTags map[enum.TagType][]tags.Tag
		err       error
	}
	type test struct {
		name string
		args args
		want want
	}

	generalTag1 := tags.Tag{
		ID:   1,
		Name: "general",
		Type: enum.TagGeneral,
	}
	generalTag2 := tags.Tag{
		ID:   2,
		Name: "general 2",
		Type: enum.TagGeneral,
	}
	peopleTag1 := tags.Tag{
		ID:   1,
		Name: "people",
		Type: enum.TagPeople,
	}
	peopleTag2 := tags.Tag{
		ID:   2,
		Name: "people",
		Type: enum.TagPeople,
	}

	mixedTags := []tags.Tag{
		generalTag1,
		peopleTag1,
	}
	generalTags := []tags.Tag{
		generalTag1,
		generalTag2,
	}
	peopleTags := []tags.Tag{
		peopleTag1,
		peopleTag2,
	}
	emptyTags := []tags.Tag{}

	mixedMap := map[enum.TagType][]tags.Tag{
		enum.TagGeneral: {
			generalTag1,
		},
		enum.TagPeople: {
			peopleTag1,
		},
	}
	generalMap := map[enum.TagType][]tags.Tag{
		enum.TagGeneral: generalTags,
	}
	peopleMap := map[enum.TagType][]tags.Tag{
		enum.TagPeople: peopleTags,
	}
	emptyMap := map[enum.TagType][]tags.Tag{}

	tests := []test{
		{
			name: "simple tag list",
			args: args{
				allTags: mixedTags,
			},
			want: want{
				splitTags: mixedMap,
			},
		},
		{
			name: "general tag list",
			args: args{
				allTags: generalTags,
			},
			want: want{
				splitTags: generalMap,
			},
		},
		{
			name: "people tag list",
			args: args{
				allTags: peopleTags,
			},
			want: want{
				splitTags: peopleMap,
			},
		},
		{
			name: "empty tag list",
			args: args{
				allTags: emptyTags,
			},
			want: want{
				splitTags: emptyMap,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := NewTagService(nil)

			splitTags, err := service.SeperateTagTypes(context.Background(), test.args.allTags)
			assert.Equal(t, test.want.err, err)
			assert.Equal(t, test.want.splitTags, splitTags)
		})
	}
}

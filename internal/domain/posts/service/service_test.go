package service

import (
	"context"
	"errors"
	"goserv/internal/domain/posts"
	"goserv/internal/domain/posts/repository"
	"goserv/internal/domain/tags"
	"goserv/internal/static/enum"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostService_GetPost(t *testing.T) {
	type args struct {
		postID int
	}
	type want struct {
		post *posts.Post
		err  error
	}
	type test struct {
		name string
		args args
		want want
	}

	basicArgs := args{
		postID: 1,
	}

	basicTags := []tags.Tag{
		{
			ID:   1,
			Name: "tag1",
			Type: enum.TagGeneral,
		},
		{
			ID:   2,
			Name: "tag2",
			Type: enum.TagPeople,
		},
	}

	basicPost := posts.Post{
		ID:        1,
		Title:     "title",
		MediaType: enum.MediaImage,
		Filename:  "filename",
		FileExt:   ".ext",
		OwnerID:   1,

		Tags: basicTags,
	}

	tests := []test{
		{
			name: "simple get post",
			args: basicArgs,
			want: want{
				post: &basicPost,
				err:  nil,
			},
		},
		{
			name: "error get post",
			args: basicArgs,
			want: want{
				post: nil,
				err:  errors.New("test error"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			postRepo := &repository.PostMock{
				GetPostFunc: func(ctx context.Context, postID int) (*posts.Post, error) {
					return test.want.post, test.want.err
				},
			}

			service := NewPostService(postRepo)

			post, err := service.GetPost(context.Background(), test.args.postID)
			assert.Equal(t, test.want.err, err)
			assert.Equal(t, test.want.post, post)
		})
	}
}

func TestPostService_ListPosts(t *testing.T) {
	type want struct {
		posts []posts.Post
		err   error
	}
	type test struct {
		name string
		want want
	}

	basicTags := []tags.Tag{
		{
			ID:   1,
			Name: "tag1",
			Type: enum.TagGeneral,
		},
		{
			ID:   2,
			Name: "tag2",
			Type: enum.TagPeople,
		},
	}

	basicPosts := []posts.Post{
		{
			ID:        1,
			Title:     "title",
			MediaType: enum.MediaImage,
			Filename:  "filename",
			FileExt:   ".ext",
			OwnerID:   1,

			Tags: basicTags,
		},
		{
			ID:        2,
			Title:     "title",
			MediaType: enum.MediaVideo,
			Filename:  "filename",
			FileExt:   ".ext",
			OwnerID:   1,

			Tags: basicTags,
		},
	}

	tests := []test{
		{
			name: "simple list posts",
			want: want{
				posts: basicPosts,
				err:   nil,
			},
		},
		{
			name: "error list posts",
			want: want{
				posts: nil,
				err:   errors.New("test error"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			postRepo := &repository.PostMock{
				ListPostsFunc: func(ctx context.Context) ([]posts.Post, error) {
					return test.want.posts, test.want.err
				},
			}

			service := NewPostService(postRepo)

			posts, err := service.ListPosts(context.Background())
			assert.Equal(t, test.want.err, err)
			assert.Equal(t, test.want.posts, posts)
		})
	}
}

func TestPostService_ListUserPosts(t *testing.T) {
	type args struct {
		userID int
	}
	type want struct {
		posts []posts.Post
		err   error
	}
	type test struct {
		name string
		args args
		want want
	}

	basicArgs := args{
		userID: 1,
	}

	basicTags := []tags.Tag{
		{
			ID:   1,
			Name: "tag1",
			Type: enum.TagGeneral,
		},
		{
			ID:   2,
			Name: "tag2",
			Type: enum.TagPeople,
		},
	}

	basicPosts := []posts.Post{
		{
			ID:        1,
			Title:     "title",
			MediaType: enum.MediaImage,
			Filename:  "filename",
			FileExt:   ".ext",
			OwnerID:   1,

			Tags: basicTags,
		},
		{
			ID:        2,
			Title:     "title",
			MediaType: enum.MediaVideo,
			Filename:  "filename",
			FileExt:   ".ext",
			OwnerID:   1,

			Tags: basicTags,
		},
	}

	tests := []test{
		{
			name: "simple list user posts",
			args: basicArgs,
			want: want{
				posts: basicPosts,
				err:   nil,
			},
		},
		{
			name: "error list user posts",
			args: basicArgs,
			want: want{
				posts: nil,
				err:   errors.New("test error"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			postRepo := &repository.PostMock{
				ListUserPostsFunc: func(ctx context.Context, userID int) ([]posts.Post, error) {
					return test.want.posts, test.want.err
				},
			}

			service := NewPostService(postRepo)

			posts, err := service.ListUserPosts(context.Background(), test.args.userID)
			assert.Equal(t, test.want.err, err)
			assert.Equal(t, test.want.posts, posts)
		})
	}
}

func TestPostService_ListUserFavs(t *testing.T) {
	type args struct {
		userID int
	}
	type want struct {
		posts []posts.Post
		err   error
	}
	type test struct {
		name string
		args args
		want want
	}

	basicArgs := args{
		userID: 1,
	}

	basicTags := []tags.Tag{
		{
			ID:   1,
			Name: "tag1",
			Type: enum.TagGeneral,
		},
		{
			ID:   2,
			Name: "tag2",
			Type: enum.TagPeople,
		},
	}

	basicPosts := []posts.Post{
		{
			ID:        1,
			Title:     "title",
			MediaType: enum.MediaImage,
			Filename:  "filename",
			FileExt:   ".ext",
			OwnerID:   1,

			Tags: basicTags,
		},
		{
			ID:        2,
			Title:     "title",
			MediaType: enum.MediaVideo,
			Filename:  "filename",
			FileExt:   ".ext",
			OwnerID:   1,

			Tags: basicTags,
		},
	}

	tests := []test{
		{
			name: "simple list user posts",
			args: basicArgs,
			want: want{
				posts: basicPosts,
				err:   nil,
			},
		},
		{
			name: "error list user posts",
			args: basicArgs,
			want: want{
				posts: nil,
				err:   errors.New("test error"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			postRepo := &repository.PostMock{
				ListUserFavsFunc: func(ctx context.Context, userID int) ([]posts.Post, error) {
					return test.want.posts, test.want.err
				},
			}

			service := NewPostService(postRepo)

			posts, err := service.ListUserFavs(context.Background(), test.args.userID)
			assert.Equal(t, test.want.err, err)
			assert.Equal(t, test.want.posts, posts)
		})
	}
}

func TestPostService_FavouritePost(t *testing.T) {
	type args struct {
		postID int
		userID int
	}
	type want struct {
		err error
	}
	type test struct {
		name string
		args args
		want want
	}

	basicArgs := args{
		postID: 1,
		userID: 1,
	}

	tests := []test{
		{
			name: "simple favourite post",
			args: basicArgs,
			want: want{
				err: nil,
			},
		},
		{
			name: "error favourite post",
			args: basicArgs,
			want: want{
				err: errors.New("test error"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			postRepo := &repository.PostMock{
				FavouritePostFunc: func(ctx context.Context, postID int, userID int) error {
					return test.want.err
				},
			}

			service := NewPostService(postRepo)

			err := service.FavouritePost(context.Background(), test.args.postID, test.args.userID)
			assert.Equal(t, test.want.err, err)
		})
	}
}

func TestPostService_UnfavouritePost(t *testing.T) {
	type args struct {
		postID int
		userID int
	}
	type want struct {
		err error
	}
	type test struct {
		name string
		args args
		want want
	}

	basicArgs := args{
		postID: 1,
		userID: 1,
	}

	tests := []test{
		{
			name: "simple favourite post",
			args: basicArgs,
			want: want{
				err: nil,
			},
		},
		{
			name: "error favourite post",
			args: basicArgs,
			want: want{
				err: errors.New("test error"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			postRepo := &repository.PostMock{
				UnfavouritePostFunc: func(ctx context.Context, postID int, userID int) error {
					return test.want.err
				},
			}

			service := NewPostService(postRepo)

			err := service.UnfavouritePost(context.Background(), test.args.postID, test.args.userID)
			assert.Equal(t, test.want.err, err)
		})
	}
}

func TestPostService_GetPostWithFavouriteStatus(t *testing.T) {
	type args struct {
		postID int
		userID int
	}
	type want struct {
		post  *posts.Post
		isFav bool
		err   error
	}
	type test struct {
		name string
		args args
		want want
	}

	basicArgs := args{
		postID: 1,
		userID: 1,
	}

	basicTags := []tags.Tag{
		{
			ID:   1,
			Name: "tag1",
			Type: enum.TagGeneral,
		},
		{
			ID:   2,
			Name: "tag2",
			Type: enum.TagPeople,
		},
	}

	basicPost := posts.Post{
		ID:        1,
		Title:     "title",
		MediaType: enum.MediaImage,
		Filename:  "filename",
		FileExt:   ".ext",
		OwnerID:   1,

		Tags: basicTags,
	}

	tests := []test{
		{
			name: "simple get post with favourite status true",
			args: basicArgs,
			want: want{
				post:  &basicPost,
				isFav: true,
				err:   nil,
			},
		},
		{
			name: "simple get post with favourite status false",
			args: basicArgs,
			want: want{
				post:  &basicPost,
				isFav: false,
				err:   nil,
			},
		},
		{
			name: "error get post with favourite status",
			args: basicArgs,
			want: want{
				post:  nil,
				isFav: false,
				err:   errors.New("test error"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			postRepo := &repository.PostMock{
				GetPostWithFavouriteStatusFunc: func(ctx context.Context, postID int, userID int) (*posts.Post, bool, error) {
					return test.want.post, test.want.isFav, test.want.err
				},
			}

			service := NewPostService(postRepo)

			post, isFav, err := service.GetPostWithFavouriteStatus(context.Background(), test.args.postID, test.args.userID)
			assert.Equal(t, test.want.err, err)
			assert.Equal(t, test.want.isFav, isFav)
			assert.Equal(t, test.want.post, post)
		})
	}
}

package repository

import (
	"context"
	"goserv/internal/domain/posts"
)

type PostMock struct {
	AddPostFunc                    func(ctx context.Context, post *posts.Post, userID int) (int, error)
	DeletePostFunc                 func(ctx context.Context, postID int) error
	GetPostFunc                    func(ctx context.Context, postID int) (*posts.Post, error)
	ListPostsFunc                  func(ctx context.Context) ([]posts.Post, error)
	ListUserPostsFunc              func(ctx context.Context, userID int) ([]posts.Post, error)
	ListUserFavsFunc               func(ctx context.Context, userID int) ([]posts.Post, error)
	FavouritePostFunc              func(ctx context.Context, postID int, userID int) error
	UnfavouritePostFunc            func(ctx context.Context, postID int, userID int) error
	GetPostWithFavouriteStatusFunc func(ctx context.Context, postID int, userID int) (*posts.Post, bool, error)
}

func (m *PostMock) AddPost(ctx context.Context, post *posts.Post, userID int) (int, error) {
	return m.AddPostFunc(ctx, post, userID)
}

func (m *PostMock) DeletePost(ctx context.Context, postID int) error {
	return m.DeletePostFunc(ctx, postID)
}

func (m *PostMock) GetPost(ctx context.Context, postID int) (*posts.Post, error) {
	return m.GetPostFunc(ctx, postID)
}

func (m *PostMock) ListPosts(ctx context.Context) ([]posts.Post, error) {
	return m.ListPostsFunc(ctx)
}

func (m *PostMock) ListUserPosts(ctx context.Context, userID int) ([]posts.Post, error) {
	return m.ListUserPostsFunc(ctx, userID)
}

func (m *PostMock) ListUserFavs(ctx context.Context, userID int) ([]posts.Post, error) {
	return m.ListUserFavsFunc(ctx, userID)
}

func (m *PostMock) FavouritePost(ctx context.Context, postID int, userID int) error {
	return m.FavouritePostFunc(ctx, postID, userID)
}

func (m *PostMock) UnfavouritePost(ctx context.Context, postID int, userID int) error {
	return m.UnfavouritePostFunc(ctx, postID, userID)
}

func (m *PostMock) GetPostWithFavouriteStatus(ctx context.Context, postID int, userID int) (*posts.Post, bool, error) {
	return m.GetPostWithFavouriteStatusFunc(ctx, postID, userID)
}

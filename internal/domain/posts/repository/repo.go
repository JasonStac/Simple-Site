package repository

import (
	"context"
	"goserv/ent/gen"
	entPost "goserv/ent/gen/post"
	"goserv/ent/gen/user"
	"goserv/internal/domain/posts"
	"goserv/internal/models"
)

type Post interface {
	AddPost(ctx context.Context, post *posts.Post, userID int) (int, error)
	DeletePost(ctx context.Context, postID int) error
	ListPosts(ctx context.Context) ([]posts.Post, error)
	ListUserPosts(ctx context.Context, userID int) ([]posts.Post, error)
	ListUserFavs(ctx context.Context, userID int) ([]posts.Post, error)
}

type postRepository struct {
	client *gen.Client
}

func NewPostRepository(client *gen.Client) *postRepository {
	return &postRepository{client: client}
}

func (repo *postRepository) AddPost(ctx context.Context, post *posts.Post, userID int) (int, error) {
	savedPost, err := repo.client.Post.Create().SetTitle(post.Title).SetMediaType(entPost.MediaType(post.MediaType)).SetFilename(post.Filename).SetOwnerID(userID).Save(ctx)
	if err != nil {
		return -1, err
	}
	return savedPost.ID, nil
}

func (repo *postRepository) DeletePost(ctx context.Context, postID int) error {
	_, err := repo.client.Post.Delete().Where(entPost.IDEQ(postID)).Exec(ctx)
	return err
}

func (repo *postRepository) ListPosts(ctx context.Context) ([]posts.Post, error) {
	entPosts, err := repo.client.Post.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	returnPosts := make([]posts.Post, len(entPosts))
	for i, p := range entPosts {
		returnPosts[i] = posts.Post{ID: p.ID, Title: p.Title, MediaType: models.MediaType(p.MediaType), Filename: p.Filename}
	}
	return returnPosts, err
}

func (repo *postRepository) ListUserPosts(ctx context.Context, userID int) ([]posts.Post, error) {
	entPosts, err := repo.client.Post.Query().Where(entPost.IDEQ(userID)).All(ctx)
	returnPosts := make([]posts.Post, len(entPosts))
	for i, p := range entPosts {
		returnPosts[i] = posts.Post{ID: p.ID, Title: p.Title, MediaType: models.MediaType(p.MediaType), Filename: p.Filename}
	}
	return returnPosts, err
}

func (repo *postRepository) ListUserFavs(ctx context.Context, userID int) ([]posts.Post, error) {
	entPosts, err := repo.client.User.Query().Where(user.IDEQ(userID)).QueryFavourites().All(ctx)
	returnPosts := make([]posts.Post, len(entPosts))
	for i, p := range entPosts {
		returnPosts[i] = posts.Post{ID: p.ID, Title: p.Title, MediaType: models.MediaType(p.MediaType), Filename: p.Filename}
	}
	return returnPosts, err
}

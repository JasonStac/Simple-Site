package repository

import (
	"context"
	"goserv/ent/gen"
	entPost "goserv/ent/gen/post"
	entUser "goserv/ent/gen/user"
	"goserv/internal/domain/artists"
	"goserv/internal/domain/posts"
	"goserv/internal/domain/tags"
	"goserv/internal/models"
	"goserv/internal/utils/errors"
)

type Post interface {
	AddPost(ctx context.Context, post *posts.Post, userID int) (int, error)
	DeletePost(ctx context.Context, postID int) error
	GetPost(ctx context.Context, postID int) (*posts.Post, error)
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
	var tagIDs []int
	for _, t := range *post.Tags {
		tagIDs = append(tagIDs, t.ID)
	}
	var artistIDs []int
	for _, a := range *post.Artists {
		artistIDs = append(artistIDs, a.ID)
	}
	// TODO: look at figuring out proper domain tag/artist into ent tag/artist and using those

	savedPost, err := repo.client.Post.
		Create().
		SetTitle(post.Title).
		SetMediaType(entPost.MediaType(post.MediaType)).
		SetFilename(post.Filename).
		SetOwnerID(userID).
		AddTagIDs(tagIDs...).
		AddArtistIDs(artistIDs...).
		Save(ctx)
	if err != nil {
		return -1, err
	}
	return savedPost.ID, nil
}

func (repo *postRepository) DeletePost(ctx context.Context, postID int) error {
	_, err := repo.client.Post.Delete().Where(entPost.IDEQ(postID)).Exec(ctx)
	return err
}

func (repo *postRepository) GetPost(ctx context.Context, postID int) (*posts.Post, error) {
	post, err := repo.client.Post.Query().Where(entPost.IDEQ(postID)).WithArtists().WithTags().First(ctx)
	if err != nil {
		if gen.IsNotFound(err) {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}

	//TODO: move conversions somewhere else
	var domainArtists []artists.Artist
	for _, a := range post.Edges.Artists {
		artist := artists.Artist{ID: a.ID, Name: a.Name}
		domainArtists = append(domainArtists, artist)
	}

	var domainTags []tags.Tag
	for _, t := range post.Edges.Tags {
		tag := tags.Tag{ID: t.ID, Name: t.Name}
		domainTags = append(domainTags, tag)
	}

	result := &posts.Post{
		ID:        post.ID,
		Title:     post.Title,
		MediaType: models.MediaType(post.MediaType),
		Filename:  post.Filename,

		Artists: &domainArtists,
		Tags:    &domainTags,
	}
	return result, nil
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
	entPosts, err := repo.client.User.Query().Where(entUser.IDEQ(userID)).QueryOwns().All(ctx)
	returnPosts := make([]posts.Post, len(entPosts))
	for i, p := range entPosts {
		returnPosts[i] = posts.Post{ID: p.ID, Title: p.Title, MediaType: models.MediaType(p.MediaType), Filename: p.Filename}
	}
	return returnPosts, err
}

func (repo *postRepository) ListUserFavs(ctx context.Context, userID int) ([]posts.Post, error) {
	entPosts, err := repo.client.User.Query().Where(entUser.IDEQ(userID)).QueryFavourites().All(ctx)
	returnPosts := make([]posts.Post, len(entPosts))
	for i, p := range entPosts {
		returnPosts[i] = posts.Post{ID: p.ID, Title: p.Title, MediaType: models.MediaType(p.MediaType), Filename: p.Filename}
	}
	return returnPosts, err
}

package repository

import (
	"context"
	"goserv/ent/gen"
	entPost "goserv/ent/gen/post"
	entUser "goserv/ent/gen/user"
	"goserv/internal/domain/posts"
	"goserv/internal/domain/tags"
	"goserv/internal/static/enum"
	"goserv/internal/utils/errors"
)

type Post interface {
	AddPost(ctx context.Context, post *posts.Post, userID int) (int, error)
	DeletePost(ctx context.Context, postID int) error
	GetPost(ctx context.Context, postID int) (*posts.Post, error)
	ListPosts(ctx context.Context) ([]posts.Post, error)
	ListUserPosts(ctx context.Context, userID int) ([]posts.Post, error)
	ListUserFavs(ctx context.Context, userID int) ([]posts.Post, error)
	FavouritePost(ctx context.Context, postID int, userID int) error
	UnfavouritePost(ctx context.Context, postID int, userID int) error
	GetPostWithFavouriteStatus(ctx context.Context, postID int, userID int) (*posts.Post, bool, error)
}

type postRepository struct {
	client *gen.Client
}

func NewPostRepository(client *gen.Client) *postRepository {
	return &postRepository{client: client}
}

func (repo *postRepository) AddPost(ctx context.Context, post *posts.Post, userID int) (int, error) {
	tagIDs := make([]int, len(post.Tags))
	for i := range post.Tags {
		tagIDs[i] = post.Tags[i].ID
	}

	savedPost, err := repo.client.Post.
		Create().
		SetTitle(post.Title).
		SetMediaType(entPost.MediaType(post.MediaType)).
		SetFilename(post.Filename).
		SetFileExt(post.FileExt).
		SetOwnerID(userID).
		AddTagIDs(tagIDs...).
		Save(ctx)
	if err != nil {
		return 0, err
	}
	return savedPost.ID, nil
}

func (repo *postRepository) DeletePost(ctx context.Context, postID int) error {
	return repo.client.Post.DeleteOneID(postID).Exec(ctx)
}

func (repo *postRepository) GetPost(ctx context.Context, postID int) (*posts.Post, error) {
	post, err := repo.client.Post.Query().Where(entPost.IDEQ(postID)).WithTags().Only(ctx)
	if err != nil {
		if gen.IsNotFound(err) {
			return nil, errors.ErrNotFound
		}
		return nil, err
	}

	//TODO: move conversions somewhere else
	domainTags := make([]tags.Tag, len(post.Edges.Tags))
	for i := range post.Edges.Tags {
		domainTags[i] = tags.Tag{
			ID:   post.Edges.Tags[i].ID,
			Type: enum.TagType(post.Edges.Tags[i].TagType),
			Name: post.Edges.Tags[i].Name,
		}
	}

	result := &posts.Post{
		ID:        post.ID,
		Title:     post.Title,
		MediaType: enum.MediaType(post.MediaType),
		Filename:  post.Filename,
		FileExt:   post.FileExt,
		OwnerID:   post.UserOwns,

		Tags: domainTags,
	}
	return result, nil
}

func (repo *postRepository) ListPosts(ctx context.Context) ([]posts.Post, error) {
	entPosts, err := repo.client.Post.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	returnPosts := make([]posts.Post, len(entPosts))
	for i := range entPosts {
		returnPosts[i] = posts.Post{
			ID:        entPosts[i].ID,
			Title:     entPosts[i].Title,
			MediaType: enum.MediaType(entPosts[i].MediaType),
			Filename:  entPosts[i].Filename,
			FileExt:   entPosts[i].FileExt,
		}
	}
	return returnPosts, err
}

func (repo *postRepository) ListUserPosts(ctx context.Context, userID int) ([]posts.Post, error) {
	entPosts, err := repo.client.User.Query().Where(entUser.IDEQ(userID)).QueryOwns().All(ctx)
	returnPosts := make([]posts.Post, len(entPosts))
	for i := range entPosts {
		returnPosts[i] = posts.Post{
			ID:        entPosts[i].ID,
			Title:     entPosts[i].Title,
			MediaType: enum.MediaType(entPosts[i].MediaType),
			Filename:  entPosts[i].Filename,
			FileExt:   entPosts[i].FileExt,
		}
	}
	return returnPosts, err
}

func (repo *postRepository) ListUserFavs(ctx context.Context, userID int) ([]posts.Post, error) {
	entPosts, err := repo.client.User.Query().Where(entUser.IDEQ(userID)).QueryFavourites().All(ctx)
	returnPosts := make([]posts.Post, len(entPosts))
	for i := range entPosts {
		returnPosts[i] = posts.Post{
			ID:        entPosts[i].ID,
			Title:     entPosts[i].Title,
			MediaType: enum.MediaType(entPosts[i].MediaType),
			Filename:  entPosts[i].Filename,
			FileExt:   entPosts[i].FileExt,
		}
	}
	return returnPosts, err
}

func (repo *postRepository) FavouritePost(ctx context.Context, postID int, userID int) error {
	return repo.client.User.UpdateOneID(userID).AddFavouriteIDs(postID).Exec(ctx)
}

func (repo *postRepository) UnfavouritePost(ctx context.Context, postID int, userID int) error {
	return repo.client.User.UpdateOneID(userID).RemoveFavouriteIDs(postID).Exec(ctx)
}

func (repo *postRepository) GetPostWithFavouriteStatus(ctx context.Context, postID int, userID int) (*posts.Post, bool, error) {
	post, err := repo.client.Post.
		Query().
		Where(entPost.IDEQ(postID)).
		WithTags().
		WithFavouritedBy(func(q *gen.UserQuery) {
			q.Where(entUser.ID(userID))
		}).
		Only(ctx)
	if err != nil {
		if gen.IsNotFound(err) {
			return nil, false, errors.ErrNotFound
		}
		return nil, false, err
	}

	domainTags := make([]tags.Tag, len(post.Edges.Tags))
	for i := range post.Edges.Tags {
		domainTags[i] = tags.Tag{
			ID:   post.Edges.Tags[i].ID,
			Type: enum.TagType(post.Edges.Tags[i].TagType),
			Name: post.Edges.Tags[i].Name,
		}
	}

	result := &posts.Post{
		ID:        post.ID,
		Title:     post.Title,
		MediaType: enum.MediaType(post.MediaType),
		Filename:  post.Filename,
		FileExt:   post.FileExt,
		OwnerID:   post.UserOwns,

		Tags: domainTags,
	}
	return result, len(post.Edges.FavouritedBy) > 0, nil
}

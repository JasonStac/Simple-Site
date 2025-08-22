package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"goserv/internal/domain/posts"
	"goserv/internal/domain/posts/repository"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

type PostService struct {
	repo repository.Post
}

func NewPostService(repo repository.Post) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) AddPost(ctx context.Context, post *posts.Post, content multipart.File, userID int) error {
	tempFile, err := os.CreateTemp("tmp", "upload-*")
	if err != nil {
		return err
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	/////move to seperate function
	hasher := sha256.New()
	multiWriter := io.MultiWriter(tempFile, hasher)

	_, err = io.Copy(multiWriter, content)
	tempFile.Close()
	if err != nil {
		return err
	}

	hashBytes := hasher.Sum(nil)
	hashHex := hex.EncodeToString(hashBytes)
	ext := strings.ToLower(filepath.Ext(post.Filename))

	dir1 := hashHex[0:2]
	dir2 := hashHex[2:4]
	finalDir := filepath.Join("content", dir1, dir2)
	finalName := hashHex + ext
	finalPath := filepath.Join(finalDir, finalName)
	/////move to seperate function

	if err := os.MkdirAll(finalDir, 0755); err != nil {
		return err
	}

	post.Filename = finalName
	postID, err := s.repo.AddPost(ctx, post, userID)
	if err != nil {
		return err
	}

	if err := os.Rename(tempFile.Name(), finalPath); err != nil {
		dbErr := s.repo.DeletePost(ctx, postID)
		if dbErr != nil {
			log.Printf("Error deleting post from db, %v\n", dbErr)
		}
		return err
	}

	return nil
}

func (s *PostService) GetPost(ctx context.Context, postID int) (*posts.Post, error) {
	return s.repo.GetPost(ctx, postID)
}

func (s *PostService) ListPosts(ctx context.Context) ([]posts.Post, error) {
	return s.repo.ListPosts(ctx)
}

func (s *PostService) ListUserPosts(ctx context.Context, userID int) ([]posts.Post, error) {
	return s.repo.ListUserPosts(ctx, userID)
}

func (s *PostService) ListUserFavs(ctx context.Context, userID int) ([]posts.Post, error) {
	return s.repo.ListUserFavs(ctx, userID)
}

func (s *PostService) DeletePost(ctx context.Context, postID int, filepath string) error {
	err := s.repo.DeletePost(ctx, postID)
	if err != nil {
		return err
	}

	err = os.Remove(filepath)
	if err != nil {
		log.Printf("Failed to remove file during delete for file: %s\n", filepath)
		return nil
	}
	return nil
}

func (s *PostService) FavouritePost(ctx context.Context, postID int, userID int) error {
	return s.repo.FavouritePost(ctx, postID, userID)
}

func (s *PostService) UnfavouritePost(ctx context.Context, postID int, userID int) error {
	return s.repo.UnfavouritePost(ctx, postID, userID)
}

func (s *PostService) GetPostWithFavourite(ctx context.Context, postID int, userID int) (*posts.Post, bool, error) {
	if userID == -1 {
		post, err := s.repo.GetPost(ctx, postID)
		return post, false, err
	}
	return s.repo.GetPostWithFavourite(ctx, postID, userID)
}

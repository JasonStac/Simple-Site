package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"goserv/internal/domain/posts"
	"goserv/internal/domain/posts/repository"
	"goserv/internal/static/constant"
	"goserv/internal/static/enum"
	"goserv/internal/utils"
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
	//TODO: add extension validation based on content type

	finalDir := filepath.Join("content", hashHex[0:2], hashHex[2:4])
	finalName := hashHex + post.Title
	finalPath := filepath.Join(finalDir, finalName+ext)

	if err := os.MkdirAll(finalDir, 0755); err != nil {
		return err
	}

	post.Filename = finalName
	post.FileExt = ext
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

	switch post.MediaType {
	case enum.MediaImage:
		if err := utils.CreateImageThumbnail(finalDir, finalName, ext); err != nil {
			s.cleanupBadAdd(ctx, postID, finalPath)
			return err
		}
	case enum.MediaVideo:
		if err := utils.ExctractVideoThumbnail(finalName, ext); err != nil {
			s.cleanupBadAdd(ctx, postID, finalPath)
			return err
		}
	case enum.MediaAudio, enum.MediaBook:
	default:
		s.cleanupBadAdd(ctx, postID, finalPath)
		return errors.New("invalid media type")
	}
	return nil
}

func (s *PostService) cleanupBadAdd(ctx context.Context, postID int, path string) {
	if dbErr := s.repo.DeletePost(ctx, postID); dbErr != nil {
		log.Printf("Error deleting post from db, %v\n", dbErr)
	}
	if osErr := os.Remove(path); osErr != nil {
		log.Printf("Error deleting file: %s, %v\n", path, osErr)
	}
}

func (s *PostService) GetPost(ctx context.Context, postID int) (*posts.Post, error) {
	return s.repo.GetPost(ctx, postID)
}

func (s *PostService) ListPosts(ctx context.Context) ([]posts.Post, error) {
	return s.repo.ListPosts(ctx)
}

func (s *PostService) ListPostsByGeneralTag(ctx context.Context, tagName string) ([]posts.Post, error) {
	return s.repo.ListPostsByGeneralTag(ctx, tagName)
}

func (s *PostService) ListPostsByPeopleTag(ctx context.Context, tagName string) ([]posts.Post, error) {
	return s.repo.ListPostsByPeopleTag(ctx, tagName)
}

func (s *PostService) ListSearchedPosts(ctx context.Context, query string) ([]posts.Post, error) {
	// If query has no colons, treat it as a simple title search
	if !strings.Contains(query, ":") {
		return s.repo.ListTitleSearchPosts(ctx, query)
	}

	// If query has single colon and no spaces, treat it as a simple single tag search
	if strings.Count(query, ":") == 1 && !strings.Contains(query, " ") {
		parts := strings.SplitN(query, ":", 2)
		tagType := parts[0]
		tagName := parts[1]

		switch tagType {
		case "tag":
			return s.ListPostsByGeneralTag(ctx, tagName)
		case "person":
			return s.ListPostsByPeopleTag(ctx, tagName)
		default:
			return nil, errors.New("invalid tag type")
		}
	}

	// TODO: implement elasticsearch query
	// Handle complex queries with multiple search fields
	return make([]posts.Post, 0), nil
}

func (s *PostService) ListUserPosts(ctx context.Context, userID int) ([]posts.Post, error) {
	return s.repo.ListUserPosts(ctx, userID)
}

func (s *PostService) ListUserFavs(ctx context.Context, userID int) ([]posts.Post, error) {
	return s.repo.ListUserFavs(ctx, userID)
}

func (s *PostService) DeletePost(ctx context.Context, postID int, filename string, fileExt string) error {
	err := s.repo.DeletePost(ctx, postID)
	if err != nil {
		return err
	}

	dataPath := filepath.Join(filename[0:2], filename[2:4], filename)
	err = os.Remove(filepath.Join("content", dataPath+fileExt))
	if err != nil {
		log.Printf("Failed to remove content file during delete for data: %s, %v\n", dataPath, err)
	}

	err = os.Remove(filepath.Join("thumbnails", dataPath+constant.ThumbnailExt))
	if err != nil {
		log.Printf("Failed to remove thumbnail file during delete for data: %s, %v\n", dataPath, err)
	}
	return nil
}

func (s *PostService) FavouritePost(ctx context.Context, postID int, userID int) error {
	return s.repo.FavouritePost(ctx, postID, userID)
}

func (s *PostService) UnfavouritePost(ctx context.Context, postID int, userID int) error {
	return s.repo.UnfavouritePost(ctx, postID, userID)
}

func (s *PostService) GetPostWithFavouriteStatus(ctx context.Context, postID int, userID int) (*posts.Post, bool, error) {
	if userID == 0 {
		post, err := s.repo.GetPost(ctx, postID)
		return post, false, err
	}
	return s.repo.GetPostWithFavouriteStatus(ctx, postID, userID)
}

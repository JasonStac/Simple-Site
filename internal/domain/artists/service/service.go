package service

import (
	"context"
	"goserv/internal/domain/artists"
	"goserv/internal/domain/artists/repository"
)

type ArtistService struct {
	repo repository.Artist
}

func NewArtistService(repo repository.Artist) *ArtistService {
	return &ArtistService{repo: repo}
}

func (s *ArtistService) ListArtists(ctx context.Context) ([]artists.Artist, error) {
	return s.repo.ListArtists(ctx)
}

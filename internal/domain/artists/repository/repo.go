package repository

// import (
// 	"context"
// 	"goserv/ent/gen"
// 	"goserv/internal/domain/artists"
// )

// type Artist interface {
// 	AddArtist(ctx context.Context, name string) (int, error)
// 	ListArtists(ctx context.Context) ([]artists.Artist, error)
// }

// type artistRepository struct {
// 	client *gen.Client
// }

// func NewArtistRepository(client *gen.Client) *artistRepository {
// 	return &artistRepository{client: client}
// }

// func (repo *artistRepository) AddArtist(ctx context.Context, name string) (int, error) {
// 	entArtist, err := repo.client.Artist.Create().SetName(name).Save(ctx)
// 	if err != nil {
// 		return -1, err
// 	}
// 	return entArtist.ID, nil
// }

// func (repo *artistRepository) ListArtists(ctx context.Context) ([]artists.Artist, error) {
// 	entArtists, err := repo.client.Artist.Query().All(ctx)
// 	returnArtists := make([]artists.Artist, len(entArtists))
// 	for i := range entArtists {
// 		returnArtists[i] = artists.Artist{ID: entArtists[i].ID, Name: entArtists[i].Name}
// 	}
// 	return returnArtists, err
// }

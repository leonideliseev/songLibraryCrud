package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leonideliseev/songLibraryCrud/models/sqlc/queries"
)

type SongsPostgres struct {
	queries *queries.Queries
}

func NewSongsPostgres(conn *pgxpool.Pool) *SongsPostgres {
    queries := queries.New(conn)

	return &SongsPostgres{
		queries: queries,
	}
}

func (d *SongsPostgres) GetAll(ctx context.Context, limit, offset int) ([]queries.Song, error) {
	songs, err := d.queries.GetSongs(ctx, queries.GetSongsParams{
        Limit: int32(limit),
        Offset: int32(offset),
    })
	if err != nil {
		return nil, err
	}

	return songs, nil
}

func (d *SongsPostgres) CreateSong(ctx context.Context, s queries.Song) (queries.Song, error) {
	song, err := d.queries.CreateSong(ctx, queries.CreateSongParams{
		GroupName: s.GroupName,
		Name: s.Name,
		ReleaseDate: s.ReleaseDate,
		Text: s.Text,
		Link: s.Link,
	})
	if err != nil {
		return queries.Song{}, err
	}

	return song, nil
}

func (d *SongsPostgres) GetSong(ctx context.Context, group, name string) (queries.Song, error) {
	song, err := d.queries.GetSong(ctx, queries.GetSongParams{
		GroupName: group,
		Name: name,
	})
	if err != nil {
		return queries.Song{}, err
	}

	return song, nil
}

func (d *SongsPostgres) DeleteSong(ctx context.Context, group, name string) error {
	return nil
}

func (d *SongsPostgres) UpdateSong(ctx context.Context, group, name string, updatedData queries.Song) (queries.Song, error) {
	return queries.Song{}, nil
}

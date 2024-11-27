package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leonideliseev/songLibraryCrud/internal/sqlc/queries"
	"github.com/leonideliseev/songLibraryCrud/internal/utils/convert/song"
	"github.com/leonideliseev/songLibraryCrud/models"
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

func (d *SongsPostgres) GetAll(ctx context.Context, limit, offset int) ([]models.Song, error) {
	songsModel, err := d.queries.GetSongs(ctx, queries.GetSongsParams{
        Limit: int32(limit),
        Offset: int32(offset),
    })
	if err != nil {
		return nil, err
	}

	songs := make([]models.Song, 0, len(songsModel))
	for _, sm := range songsModel {
		songs = append(songs, songConvert.FromModelToApp(sm))
	}

	return songs, nil
}

func (d *SongsPostgres) CreateSong(ctx context.Context, s models.Song) (models.Song, error) {
	createSong := songConvert.FromAppToModel(s)

	songModel, err := d.queries.CreateSong(ctx, queries.CreateSongParams{
		GroupName: createSong.GroupName,
		Name: createSong.Name,
		ReleaseDate: createSong.ReleaseDate,
		Text: createSong.Text,
		Link: createSong.Link,
	})
	if err != nil {
		return models.Song{}, err
	}

	song := songConvert.FromModelToApp(songModel)

	return song, nil
}

func (d *SongsPostgres) GetSong(ctx context.Context, id uuid.UUID) (models.Song, error) {
	uuid := pgtype.UUID{
		Bytes: id,
		Valid: true,
	}

	songModel, err := d.queries.GetSong(ctx, uuid)
	if err != nil {
		return models.Song{}, err
	}

	song := songConvert.FromModelToApp(songModel)

	return song, nil
}

func (d *SongsPostgres) DeleteSong(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (d *SongsPostgres) UpdateSong(ctx context.Context, id uuid.UUID, updatedData models.Song) (models.Song, error) {
	return models.Song{}, nil
}

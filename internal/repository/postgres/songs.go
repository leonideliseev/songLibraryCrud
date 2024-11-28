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

func (d *SongsPostgres) GetAll(ctx context.Context, limit, offset int) ([]*models.Song, error) {
	songsQuery, err := d.queries.GetSongs(ctx, queries.GetSongsParams{
        Limit: int32(limit),
        Offset: int32(offset),
    })
	if err != nil {
		return nil, err
	}

	songs := make([]*models.Song, 0, len(songsQuery))
	for _, sq := range songsQuery {
		songs = append(songs, songConvert.FromQueryToApp(&sq))
	}

	return songs, nil
}

func (d *SongsPostgres) CreateSong(ctx context.Context, s *models.Song) (*models.Song, error) {
	createSong := songConvert.FromAppToQuery(s)

	songQuery, err := d.queries.CreateSong(ctx, queries.CreateSongParams{
		GroupName: createSong.GroupName,
		Name: createSong.Name,
		ReleaseDate: createSong.ReleaseDate,
		Text: createSong.Text,
		Link: createSong.Link,
	})
	if err != nil {
		return nil, err
	}

	song := songConvert.FromQueryToApp(&songQuery)

	return song, nil
}

func (d *SongsPostgres) GetSong(ctx context.Context, id uuid.UUID) (*models.Song, error) {
	uuid := toUUID(id)

	songQuery, err := d.queries.GetSong(ctx, uuid)
	if err != nil {
		return nil, err
	}

	song := songConvert.FromQueryToApp(&songQuery)

	return song, nil
}

func (d *SongsPostgres) UpdateSong(ctx context.Context, updatedData *models.Song) (*models.Song, error) {
	updateSong := songConvert.FromAppToQuery(updatedData)

	songQuery, err := d.queries.UpdateSong(ctx, queries.UpdateSongParams{
		ID: updateSong.ID,
		GroupName: updateSong.GroupName,
		Name: updateSong.Name,
		ReleaseDate: updateSong.ReleaseDate,
		Text: updateSong.Text,
		Link: updateSong.Link,
	})
	if err != nil {
		return nil, err
	}

	song := songConvert.FromQueryToApp(&songQuery)

	return song, nil
}

func (d *SongsPostgres) DeleteSong(ctx context.Context, id uuid.UUID) error {
	uuid := toUUID(id)

	return d.queries.DeleteSong(ctx, uuid)
}

func toUUID(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: id,
		Valid: true,
	}
}

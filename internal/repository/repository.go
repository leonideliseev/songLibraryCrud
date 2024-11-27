package repository

import (
	"context"

	"github.com/leonideliseev/songLibraryCrud/models/sqlc/queries"
)

type Songs interface {
	GetAll(ctx context.Context, limit, offest int) ([]queries.Song, error)
	CreateSong(ctx context.Context, s queries.Song) (queries.Song, error)
	GetSong(ctx context.Context, group, name string) (queries.Song, error)
	DeleteSong(ctx context.Context, group, name string) error
	UpdateSong(ctx context.Context, group, name string, updatedData queries.Song) (queries.Song, error)
}

type Repository struct {
	Songs
}
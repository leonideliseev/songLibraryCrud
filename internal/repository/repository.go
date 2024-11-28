package repository

import (
	"context"

	"github.com/leonideliseev/songLibraryCrud/models"
    "github.com/google/uuid"
)

type Songs interface {
	GetAll(ctx context.Context, limit, offest int) ([]*models.Song, error)
	CreateSong(ctx context.Context, s *models.Song) (*models.Song, error)
	GetSong(ctx context.Context, id uuid.UUID) (*models.Song, error)
	DeleteSong(ctx context.Context, id uuid.UUID) error
	UpdateSong(ctx context.Context, updatedData *models.Song) (*models.Song, error)
}

type Repository struct {
	Songs
}
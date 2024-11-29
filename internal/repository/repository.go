package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leonideliseev/songLibraryCrud/models"
	"github.com/leonideliseev/songLibraryCrud/pkg/logging"
	"github.com/leonideliseev/songLibraryCrud/internal/repository/postgres"
)

type Songs interface {
	GetAll(ctx context.Context, limit, offest int, pagModel *models.Song) ([]models.Song, error)
	CreateSong(ctx context.Context, s *models.Song) (*models.Song, error)
	GetSong(ctx context.Context, id uuid.UUID) (*models.Song, error)
	DeleteSong(ctx context.Context, id uuid.UUID) error
	UpdateSong(ctx context.Context, updatedData *models.Song) (*models.Song, error)
}

type Repository struct {
	Songs
}

func New(db *pgxpool.Pool, log *logging.Logger) *Repository {
	return &Repository{
		Songs: postgres.NewSongsPostgres(db, log),
	}
}

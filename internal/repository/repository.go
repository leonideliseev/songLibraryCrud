package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/leonideliseev/songLibraryCrud/internal/repository/postgres"
	"github.com/leonideliseev/songLibraryCrud/models"
	"github.com/leonideliseev/songLibraryCrud/pkg/logging"
	"github.com/leonideliseev/songLibraryCrud/pkg/postgresql"
)

type Songs interface {
	GetAll(ctx context.Context, limit, offest int, pagModel *models.Song) ([]models.Song, error)
	Create(ctx context.Context, s *models.Song) (*models.Song, error)
	GetById(ctx context.Context, id uuid.UUID) (*models.Song, error)
	DeleteById(ctx context.Context, id uuid.UUID) error
	UpdateById(ctx context.Context, s *models.Song) (*models.Song, error)
}

type Repository struct {
	Songs
}

func New(db postgresql.Conn, log *logging.Logger) *Repository {
	log.Info("init repository...")
	return &Repository{
		Songs: postgres.NewSongsPostgres(db, log),
	}
}

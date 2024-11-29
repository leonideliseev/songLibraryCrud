package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/leonideliseev/songLibraryCrud/internal/repository"
	"github.com/leonideliseev/songLibraryCrud/models"
	"github.com/leonideliseev/songLibraryCrud/pkg/logging"
)

type Songs interface {
	GetAll(ctx context.Context, limit, offest int, pagModel *models.Song) ([]models.Song, error)
	Create(ctx context.Context, song *models.Song) (*models.Song, error)
	GetById(ctx context.Context, id uuid.UUID) (*models.Song, error)
	DeleteById(ctx context.Context, id uuid.UUID) error
	UpdateById(ctx context.Context, id uuid.UUID, updatedData *models.Song) (*models.Song, error)
}

type Service struct {
	Songs
}

func NewService(repos *repository.Repository, log *logging.Logger) *Service {
	return &Service{
		Songs: NewSongsService(repos.Songs, log),
	}
}

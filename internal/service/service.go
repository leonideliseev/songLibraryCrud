package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/leonideliseev/songLibraryCrud/internal/repository"
	"github.com/leonideliseev/songLibraryCrud/models"
)

type Songs interface {
	GetAll(ctx context.Context, limit, offest int) ([]*models.Song, error)
	CreateSong(ctx context.Context, song *models.Song) (*models.Song, error)
	GetSong(ctx context.Context, id uuid.UUID) (*models.Song, error)
	DeleteSong(ctx context.Context, id uuid.UUID) error
	UpdateSong(ctx context.Context, id uuid.UUID, updatedData *models.Song) (*models.Song, error)
}

type Service struct {
	Songs
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Songs: NewSongsService(repos.Songs),
	}
}

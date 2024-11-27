package service

import (
	"github.com/google/uuid"
	"github.com/leonideliseev/songLibraryCrud/internal/repository"
	"github.com/leonideliseev/songLibraryCrud/models"
)

type Songs interface {
	GetAll(limit, offest int) ([]models.Song, error)
	CreateSong(song models.Song) (models.Song, error)
	GetSong(id uuid.UUID) (models.Song, error)
	DeleteSong(id uuid.UUID) error
	UpdateSong(id uuid.UUID, updatedData models.Song) (models.Song, error)
}

type Service struct {
	Songs
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Songs: NewSongsService(repos.Songs),
	}
}

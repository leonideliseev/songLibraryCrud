package service

import (
	"github.com/leonideliseev/songLibraryCrud/internal/repository"
	"github.com/leonideliseev/songLibraryCrud/models"
)

type Songs interface {
	GetAll(limit, offest int) ([]models.Song, error)
	CreateSong(song models.Song) (models.Song, error)
	GetSong(group, name string) (models.Song, error)
	DeleteSong(group, name string) error
	UpdateSong(group, name string, updatedData models.Song) (models.Song, error)
}

type Service struct {
	Songs
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Songs: NewSongsService(repos.Songs),
	}
}

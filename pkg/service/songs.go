package service

import (
	"github.com/leonideliseev/songLibraryCrud/models"
	"github.com/leonideliseev/songLibraryCrud/pkg/repository"
)

type SongsService struct {
	repo repository.Songs
}

func NewSongsService(repo repository.Songs) *SongsService {
	return &SongsService{
		repo: repo,
	}
}

func (s *SongsService) GetAll(limit, offset int) ([]*models.Song, error) {
	return s.repo.GetAll(limit, offset)
}

func (s *SongsService) CreateSong(song models.Song) (*models.Song, error) {
	return &song, nil
}

func (s *SongsService) GetSong(group, name string) (*models.Song, error) {
	return s.repo.GetSong(group, name)
}

func (s *SongsService) DeleteSong(group, name string) error {
	return s.repo.DeleteSong(group, name)
}

func (s *SongsService) UpdateSong(group, name string, updatedData *models.Song) (*models.Song, error) {
	return s.repo.UpdateSong(group, name, updatedData)
}
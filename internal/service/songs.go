package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/leonideliseev/songLibraryCrud/internal/repository"
	"github.com/leonideliseev/songLibraryCrud/models"
)

type SongsService struct {
	repo repository.Songs
}

func NewSongsService(repo repository.Songs) *SongsService {
	return &SongsService{
		repo: repo,
	}
}

var ctx context.Context // TODO: временно

func (s *SongsService) GetAll(limit, offset int) ([]models.Song, error) {
	songs, err := s.repo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return songs, nil
}

func (s *SongsService) CreateSong(song models.Song) (models.Song, error) {
	song, err := s.repo.CreateSong(ctx, song)
	if err != nil {
		return models.Song{}, err
	}
	return song, nil
}

func (s *SongsService) GetSong(id uuid.UUID) (models.Song, error) {
	song, err := s.repo.GetSong(ctx, id)
	if err != nil {
		return models.Song{}, err
	}

	return song, nil
}

func (s *SongsService) DeleteSong(id uuid.UUID) error {
	return s.repo.DeleteSong(ctx, id)
}

func (s *SongsService) UpdateSong(id uuid.UUID, updatedData models.Song) (models.Song, error) {
	song, err := s.repo.UpdateSong(ctx, id, updatedData)
	if err != nil {
		return models.Song{}, err
	}

	return song, nil
}

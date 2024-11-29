package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/leonideliseev/songLibraryCrud/internal/repository"
	songConvert "github.com/leonideliseev/songLibraryCrud/internal/utils/convert/song"
	"github.com/leonideliseev/songLibraryCrud/models"
	"github.com/leonideliseev/songLibraryCrud/pkg/logging"
)

type SongsService struct {
	log *logging.Logger
	repo repository.Songs
}

func NewSongsService(repo repository.Songs, log *logging.Logger) *SongsService {
	return &SongsService{
		log: log,
		repo: repo,
	}
}

func (s *SongsService) GetAll(ctx context.Context, limit, offset int, pagModel *models.Song) ([]models.Song, error) {
	songs, err := s.repo.GetAll(ctx, limit, offset, pagModel)
	if err != nil {
		return nil, err
	}

	return songs, nil
}

func (s *SongsService) CreateSong(ctx context.Context, song *models.Song) (*models.Song, error) {
	song, err := s.repo.CreateSong(ctx, song)
	if err != nil {
		return nil, err
	}
	return song, nil
}

func (s *SongsService) GetSong(ctx context.Context, id uuid.UUID) (*models.Song, error) {
	song, err := s.repo.GetSong(ctx, id)
	if err != nil {
		return nil, err
	}

	return song, nil
}

func (s *SongsService) UpdateSong(ctx context.Context, id uuid.UUID, updatedData *models.Song) (*models.Song, error) {
	song, err := s.repo.GetSong(ctx, id)
	if err != nil {
		return nil, err
	}

	songConvert.UniteModel(song, updatedData)

	song, err = s.repo.UpdateSong(ctx, song)
	if err != nil {
		return nil, err
	}

	return song, nil
}

func (s *SongsService) DeleteSong(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteSong(ctx, id)
}

package service

import (
	"context"

	"github.com/leonideliseev/songLibraryCrud/internal/repository"
	"github.com/leonideliseev/songLibraryCrud/internal/utils/convert/song"
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
	repoSongs, err := s.repo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	songs := make([]models.Song, 0, len(repoSongs))
	for _, rs := range repoSongs {
		songs = append(songs, songConvert.FromRepoToServ(rs))
	}

	return songs, nil
}

func (s *SongsService) CreateSong(song models.Song) (models.Song, error) {
	repoSong, err := s.repo.CreateSong(ctx, songConvert.FromServToRepo(song))
	if err != nil {
		return models.Song{}, err
	}

	song = songConvert.FromRepoToServ(repoSong)

	return song, nil
}

func (s *SongsService) GetSong(group, name string) (models.Song, error) {
	repoSong, err := s.repo.GetSong(ctx, group, name)
	if err != nil {
		return models.Song{}, err
	}

	return songConvert.FromRepoToServ(repoSong), nil
}

func (s *SongsService) DeleteSong(group, name string) error {
	return s.repo.DeleteSong(ctx, group, name)
}

func (s *SongsService) UpdateSong(group, name string, updatedData models.Song) (models.Song, error) {
	repoUpdateData := songConvert.FromServToRepo(updatedData)
	repoSong, err := s.repo.UpdateSong(ctx, group, name, repoUpdateData)
	if err != nil {
		return models.Song{}, err
	}

	return songConvert.FromRepoToServ(repoSong), nil
}

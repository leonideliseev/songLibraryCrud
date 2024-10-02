package service

import "github.com/leonideliseev/songLibraryCrud/pkg/repository"

type SongsService struct {
	repo repository.Songs
}

func NewSongsService(repo repository.Songs) *SongsService {
	return &SongsService{
		repo: repo,
	}
}
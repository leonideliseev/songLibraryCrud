package service

import "github.com/leonideliseev/songLibraryCrud/pkg/repository"

type Songs interface {
}

type Service struct {
	Songs
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Songs: NewSongsService(repos.Songs),
	}
}
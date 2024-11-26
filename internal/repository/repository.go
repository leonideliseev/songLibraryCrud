package repository

import "github.com/leonideliseev/songLibraryCrud/models"

type Songs interface {
	GetAll(limit, offest int) ([]*models.Song, error)
	CreateSong(s models.Song) (*models.Song, error)
	GetSong(group, name string) (*models.Song, error)
	DeleteSong(group, name string) error
	UpdateSong(group, name string, updatedData *models.Song) (*models.Song, error)
}

type Repository struct {
	Songs
}
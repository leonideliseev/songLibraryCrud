package postgres

import (
	"github.com/jinzhu/gorm"
)

type SongsPostgres struct {
	db *gorm.DB
}

func NewSongsPostgres(db *gorm.DB) *SongsPostgres {
	return &SongsPostgres{
		db: db,
	}
}
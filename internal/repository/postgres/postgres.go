package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leonideliseev/songLibraryCrud/internal/repository"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresRepository(db *pgxpool.Pool) *repository.Repository {
	return &repository.Repository{
		Songs: NewSongsPostgres(db),
	}
}

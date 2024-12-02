package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/leonideliseev/songLibraryCrud/internal/repository/postgres"
	"github.com/leonideliseev/songLibraryCrud/models"
	"github.com/leonideliseev/songLibraryCrud/pkg/logging"
	"github.com/leonideliseev/songLibraryCrud/pkg/postgresql"
	"github.com/spf13/viper"
)

type Songs interface {
	GetAll(ctx context.Context, limit, offest int, pagModel *models.Song) ([]models.Song, error)
	Create(ctx context.Context, s *models.Song) (*models.Song, error)
	GetById(ctx context.Context, id uuid.UUID) (*models.Song, error)
	DeleteById(ctx context.Context, id uuid.UUID) error
	UpdateById(ctx context.Context, s *models.Song) (*models.Song, error)
}

type Repository struct {
	Songs
}

func New(conn postgresql.Conn, log *logging.Logger) *Repository {
	engine := viper.GetString("repo_implement.engine")
	driver := viper.GetString("repo_implement.sqldriver")

	switch engine {
	case "postgresql", "PostgreSQL":
		switch driver {
		case "pgx/v5", "pgx":
			return &Repository{
				Songs: postgres.NewSongsPostgres(conn, log),
			}
		default:
			log.Fatalf("This driver not implemented <%s>", driver)
		}
	default:
		log.Fatalf("This data base not implemented <%s>", engine)
	}

	return nil
}

package utils

import (
	"github.com/leonideliseev/songLibraryCrud/internal/repository"
	"github.com/leonideliseev/songLibraryCrud/pkg/logging"
	"github.com/leonideliseev/songLibraryCrud/pkg/postgresql"
	"github.com/spf13/viper"
)

func RepoChoice(repo **repository.Repository, conn postgresql.Conn, log *logging.Logger) {
	engine := viper.GetString("repo_implement.engine")
	driver := viper.GetString("repo_implement.sqldriver")

	switch engine {
	case "postgresql", "PostgreSQL":
		switch driver {
		case "pgx/v5", "pgx":
			*repo = repository.New(conn, log)
		default:
			log.Fatalf("This driver not implemented <%s>", driver)
		}
	default:
		log.Fatalf("This data base not implemented <%s>", engine)
	}
}

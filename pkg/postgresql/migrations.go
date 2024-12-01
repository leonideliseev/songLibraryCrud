package postgresql

import (
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/leonideliseev/songLibraryCrud/pkg/logging"
)

func Migrate(log *logging.Logger, fs *embed.FS, cfg *Config) {
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SSLMode,
	)

	source, err := iofs.New(fs, "migrations")
	if err != nil {
		log.WithError(err).Fatal("failed to read migrations source")
	}

	instance, err := migrate.NewWithSourceInstance("iofs", source, dbUrl)
	if err != nil {
		log.WithError(err).Fatal("failed to initialization the migrations instance")
	}

	err = instance.Up()

	switch err {
	case nil:
		log.Debug("the migration schema successfully upgraded!")
	case migrate.ErrNoChange:
		log.Debug("the migration schema not changed")
	default:
		log.WithError(err).Fatal("could not apply the migration schema")
	}
}

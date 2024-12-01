package utils

import (
	"context"
	"fmt"
	"os"
	"time"

	// "github.com/golang-migrate/migrate"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leonideliseev/songLibraryCrud/internal/repository"
	"github.com/leonideliseev/songLibraryCrud/pkg/logging"
	"github.com/leonideliseev/songLibraryCrud/pkg/postgresql"
	"github.com/leonideliseev/songLibraryCrud/schema"
	"github.com/spf13/viper"
)

const (
	postgresDataBase = "postgres"
)

func PostgresPgx(cfg postgresql.Config, log *logging.Logger) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	dsn := PostgresConn(cfg)

	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.WithError(err).Fatalf("failed conn pgxpool")
	}

	return db, db.Ping(ctx)
}

func PostgresConn(cfg postgresql.Config) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", 
    cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)
}

func RepoChoice(repo *repository.Repository, log *logging.Logger) {
	engine := viper.GetString("repo_implement.engine")
	driver := viper.GetString("repo_implement.sqldriver")

	switch engine {
	case "postgresql", "PostgreSQL":
		config := postgresql.Config{
			Host:     viper.GetString("db.host"),
			Port:     viper.GetString("db.port"),
			Username: viper.GetString("db.username"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   viper.GetString("db.dbname"),
			SSLMode:  viper.GetString("db.sslmode"),
		}

		switch driver {
		case "pgx/v5", "pgx":

			testConfig := config
			// чтобы проверить наличие датабазы ставлю дефолтную постгреса
			testConfig.DBName = postgresDataBase

			conn, err := PostgresPgx(testConfig, log)
			if err != nil {
				log.WithError(err).Fatalf("failed conn to db")
			}
			defer conn.Close()

			createDatabaseIfNotExists(conn, viper.GetString("db.dbname"), log)

			conn, err = PostgresPgx(config, log)
			if err != nil {
				log.WithError(err).Fatalf("failed conn to db")
			}
			defer conn.Close()

			postgresql.Migrate(log, &schema.DB, &config)

			*repo = *repository.New(conn, log)
		default:
				log.Fatalf("This driver not implemented <%s>", driver)
		}
	default:
		log.Fatalf("This data base not implemented <%s>", engine)
	}
}

func createDatabaseIfNotExists(conn *pgxpool.Pool, dbName string, log *logging.Logger) {
	var exists bool
	err := conn.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname=$1)", dbName).Scan(&exists)
	if err != nil {
		log.WithError(err).Fatalf("failed to check if database exists")
	}

	if !exists {
		_, err := conn.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			log.WithError(err).Fatalf("failed to create database")
		}

		log.Infof("Database %s created successfully", dbName)
	}
}

/*func runMigrations(connStr string, log *logging.Logger) {
	m, err := migrate.New(
		"file://migrations",
		connStr,
	)
	if err != nil {
		log.WithError(err).Errorf("failed to create migrate instance")
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.WithError(err).Errorf("failed to apply migrations")
	}

	log.Info("Migrations applied successfully")
}*/

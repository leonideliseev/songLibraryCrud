package utils

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leonideliseev/songLibraryCrud/internal/repository"
	"github.com/leonideliseev/songLibraryCrud/internal/repository/postgres"
	"github.com/spf13/viper"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func PostgresPgx(cfg Config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	dsn := PostgresConn(cfg)

	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return db, db.Ping(ctx)
}

func PostgresConn(cfg Config) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", 
    cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)
}

func RepoChoice(repo *repository.Repository) {
	switch viper.GetString("repo_implement.engine") {
	case "postgresql":
		switch viper.GetString("repo_implement.sqldriver") {
		case "pgx/v5":
			conn, err := PostgresPgx(Config{
				Host:     viper.GetString("db.host"),
				Port:     viper.GetString("db.port"),
				Username: viper.GetString("db.username"),
				Password: os.Getenv("DB_PASSWORD"),
				DBName:   viper.GetString("db.dbname"),
				SSLMode:  viper.GetString("db.sslmode"),
			})
			if err != nil {
				//logrus.Fatalf("failed conn to db: %s", err.Error())
			}
			defer conn.Close()
		
			repo = postgres.NewPostgresRepository(conn)
		default:
			panic("")
		}
	default:
		panic("")
	}
}

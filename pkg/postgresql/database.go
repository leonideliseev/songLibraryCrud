package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leonideliseev/songLibraryCrud/pkg/logging"
)

func CreateDatabaseIfNotExists(conn *pgxpool.Pool, dbName string, log *logging.Logger) {
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

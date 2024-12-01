package postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnWithPgxPool(cfg Config) (*pgxpool.Pool, error) {
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

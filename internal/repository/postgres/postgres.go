package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leonideliseev/songLibraryCrud/internal/repository"
)

const (
	usersTable      = "users"
	todoListTable   = "todo_lists"
	usersListTable  = "user_lists"
	todoItemsTable  = "todo_items"
	listsItemsTable = "lists_items"
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

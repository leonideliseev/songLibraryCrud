package postgres

import (
	"github.com/leonideliseev/songLibraryCrud/internal/repository"
	"gorm.io/gorm"
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

func NewPostgresRepository(db *gorm.DB) *repository.Repository {
	return &repository.Repository{
		Songs: NewSongsPostgres(db),
	}
}

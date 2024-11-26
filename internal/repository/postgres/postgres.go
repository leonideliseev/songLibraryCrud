package postgres

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/leonideliseev/songLibraryCrud/internal/repository"
	"github.com/leonideliseev/songLibraryCrud/models"
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

func NewPostgresRepository(cfg Config) (*repository.Repository, error) {
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode))

	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.Song{})

	repo := &repository.Repository{
		Songs: NewSongsPostgres(db),
	}

	return repo, nil
}

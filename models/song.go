package models

import (
	"time"

	"github.com/google/uuid"
)

type Song struct {
	ID          uuid.UUID `json:"id" db:"id"`
	GroupName   string    `json:"group" db:"group_name"`
	Name        string    `json:"name" db:"name"`
	ReleaseDate time.Time `json:"release_date" db:"release_date"`
	Text        string    `json:"text" db:"text"`
	Link        string    `json:"link" db:"link"`
}

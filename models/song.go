package models

import (
	"time"

	"github.com/google/uuid"
)

type Song struct {
	ID          uuid.UUID `db:"id"`
	GroupName   string    `db:"group_name"`
	Name        string    `db:"name"`
	ReleaseDate time.Time `db:"release_date"`
	Text        string    `db:"text"`
	Link        string    `db:"link"`
}

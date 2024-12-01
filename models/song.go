package models

import "github.com/google/uuid"

type Song struct {
	ID          uuid.UUID
	GroupName   string
	Name        string
	ReleaseDate string
	Text        string
	Link        string
}

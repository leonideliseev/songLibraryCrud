package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/leonideliseev/songLibraryCrud/models"
)

func uniteModel(base, update *models.Song) {
	uniteFieldArray(&base.ID, &update.ID)
	uniteFieldString(&base.GroupName, &update.GroupName)
	uniteFieldString(&base.Name, &update.Name)
	uniteFieldTime(&base.ReleaseDate, &update.ReleaseDate)
	uniteFieldString(&base.Text, &update.Text)
	uniteFieldString(&base.Link, &update.Link)
}

func uniteFieldString(base, update *string) {
	if *update == "" {
		return
	}

	*base = *update
}

func uniteFieldArray(base, update *uuid.UUID) {
	if *update == [16]byte{} {
		return
	}

	*base = *update
}

func uniteFieldTime(base, update *time.Time) {
	if update.IsZero() {
		return
	}

	*base = *update
}
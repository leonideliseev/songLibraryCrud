package songConvert

import (
	"time"

	"github.com/google/uuid"
	"github.com/leonideliseev/songLibraryCrud/internal/handler/dto"
	"github.com/leonideliseev/songLibraryCrud/models"
)

func FromInputToModel(s *dto.RequestCreateSong, sd *dto.SongDetail) (*models.Song, error) {
	t, err := time.Parse("2006-01-02", sd.ReleaseDate)
	if err != nil {
		return nil, err
	}

	song := &models.Song{
		GroupName: s.Group,
    	Name: s.Name,
    	ReleaseDate: t,
    	Text: sd.Text,
		Link: sd.Link,
	}

	return song, nil
}

func FromInputUpdateToModel(s *dto.RequestUpdateSong) (*models.Song, error) {
	t, err := time.Parse("2006-01-02", *s.ReleaseDate)
	if err != nil {
		return nil, err
	}

	song := &models.Song{
		GroupName: fromPointerToString(s.Group),
    	Name: fromPointerToString(s.Name),
    	ReleaseDate: t,
    	Text: fromPointerToString(s.Text),
		Link: fromPointerToString(s.Link),
	}

	return song, nil
}

func fromPointerToString(ptr *string) string {
	if ptr == nil {
		return ""
	}

	return *ptr
}

func UniteModel(base, update *models.Song) {
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

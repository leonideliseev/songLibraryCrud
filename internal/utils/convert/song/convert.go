package songConvert

import (
	"github.com/google/uuid"
	"github.com/leonideliseev/songLibraryCrud/internal/handler/dto"
	"github.com/leonideliseev/songLibraryCrud/models"
)

func FromInputToModel(s *dto.RequestCreateSong, sd *dto.SongDetail) *models.Song {
	return &models.Song{
		GroupName: s.Group,
    	Name: s.Name,
    	ReleaseDate: sd.ReleaseDate,
    	Text: sd.Text,
		Link: sd.Link,
	}
}

func FromInputUpdateToModel(s *dto.RequestUpdateSong) *models.Song {
	return &models.Song{
		GroupName: fromPointerToString(s.Group),
    	Name: fromPointerToString(s.Name),
    	ReleaseDate: fromPointerToString(s.ReleaseDate),
    	Text: fromPointerToString(s.Text),
		Link: fromPointerToString(s.Link),
	}
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
	uniteFieldString(&base.ReleaseDate, &update.ReleaseDate)
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

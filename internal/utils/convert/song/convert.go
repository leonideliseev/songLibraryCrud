package songConvert

import (
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
	uniteField(&base.ID, &update.ID)
	uniteField(&base.GroupName, &update.GroupName)
	uniteField(&base.Name, &update.Name)
	uniteField(&base.ReleaseDate, &update.ReleaseDate)
	uniteField(&base.Text, &update.Text)
	uniteField(&base.Link, &update.Link)
}

func uniteField(base, update *string) {
	if *update == "" {
		return
	}

	*base = *update
}

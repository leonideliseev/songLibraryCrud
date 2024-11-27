package songConvert

import (
	"github.com/leonideliseev/songLibraryCrud/internal/handler/dto"
	"github.com/leonideliseev/songLibraryCrud/models"

	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/leonideliseev/songLibraryCrud/internal/sqlc/queries"
)

func FromInputToModel(s dto.RequestCreateSong, sd dto.SongDetail) models.Song {
	return models.Song{
		GroupName: s.Group,
    	Name: s.Name,
    	ReleaseDate: sd.ReleaseDate,
    	Text: sd.Text,
		Link: sd.Link,
	}
}

func FromInputUpdateToModel(s dto.RequestUpdateSong) models.Song {
	var song models.Song

	song.GroupName = fromPointerToString(s.Group)
	song.Name = fromPointerToString(s.Name)
	song.ReleaseDate = fromPointerToString(s.ReleaseDate)
	song.Text = fromPointerToString(s.Text)
	song.Link = fromPointerToString(s.Link)

	return song
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

func FromAppToQuery(s models.Song) queries.Song {
	parsedDate, _ := time.Parse("02.01.2006", s.ReleaseDate)

	return queries.Song{
		GroupName: s.GroupName,
		Name: s.Name,
		ReleaseDate: pgtype.Date{Time: parsedDate, Valid: true},
		Text: pgtype.Text{String: s.Text, Valid: true},
		Link: pgtype.Text{String: s.Link, Valid: true},
	}
}

func FromQueryToApp(s queries.Song) models.Song {
	return models.Song{
		ID: string(s.ID.Bytes[:]),
		GroupName: s.GroupName,
		Name: s.Name,
		ReleaseDate: s.ReleaseDate.Time.String(),
		Text: s.Text.String,
		Link: s.Link.String,
	}
}

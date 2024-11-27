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

func FromAppToModel(s models.Song) queries.Song {
	parsedDate, _ := time.Parse("02.01.2006", s.ReleaseDate)

	return queries.Song{
		GroupName: s.GroupName,
		Name: s.Name,
		ReleaseDate: pgtype.Date{Time: parsedDate, Valid: true},
		Text: pgtype.Text{String: s.Text, Valid: true},
		Link: pgtype.Text{String: s.Link, Valid: true},
	}
}

func FromModelToApp(s queries.Song) models.Song {
	return models.Song{
		ID: string(s.ID.Bytes[:]),
		GroupName: s.GroupName,
		Name: s.Name,
		ReleaseDate: s.ReleaseDate.Time.String(),
		Text: s.Text.String,
		Link: s.Link.String,
	}
}

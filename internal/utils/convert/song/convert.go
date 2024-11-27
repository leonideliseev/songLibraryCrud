package songConvert

import (
	"github.com/leonideliseev/songLibraryCrud/internal/handler/dto"
	"github.com/leonideliseev/songLibraryCrud/models"

	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/leonideliseev/songLibraryCrud/models/sqlc/queries"
)

func FromHandToServ(s dto.RequestCreateSong) models.Song {
	return models.Song{
		GroupName: s.Group,
    	Name: s.Name,
    	ReleaseDate: s.ReleaseDate,
    	Text: s.Text,
		Link: s.Link,
	}
}

func FromServToRepo(s models.Song) queries.Song {
	parsedDate, _ := time.Parse("02.01.2006", s.ReleaseDate)

	return queries.Song{
		GroupName: s.GroupName,
		Name: s.Name,
		ReleaseDate: pgtype.Date{Time: parsedDate, Valid: true},
		Text: pgtype.Text{String: s.Text, Valid: true},
		Link: pgtype.Text{String: s.Link, Valid: true},
	}
}

func FromRepoToServ(s queries.Song) models.Song {
	return models.Song{
		ID: string(s.ID.Bytes[:]),
		GroupName: s.GroupName,
		Name: s.Name,
		ReleaseDate: s.ReleaseDate.Time.String(),
		Text: s.Text.String,
		Link: s.Link.String,
	}
}

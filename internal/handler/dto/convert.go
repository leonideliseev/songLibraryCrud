package dto

import (
	"time"

	"github.com/leonideliseev/songLibraryCrud/models"
)

func FromInputToModel(s *RequestCreateSong, sd *SongDetail) (*models.Song, error) {
	t, err := time.Parse("2006-01-02", sd.ReleaseDate)
	if err != nil {
		return nil, err
	}

	song := &models.Song{
		GroupName:   s.Group,
		Name:        s.Name,
		ReleaseDate: t,
		Text:        sd.Text,
		Link:        sd.Link,
	}

	return song, nil
}

func FromInputUpdateToModel(s *RequestUpdateSong) (*models.Song, error) {
	t, err := fromPointerToDate(s.ReleaseDate)
	if err != nil {
		return nil, err
	}

	song := &models.Song{
		GroupName:   fromPointerToString(s.Group),
		Name:        fromPointerToString(s.Name),
		ReleaseDate: t,
		Text:        fromPointerToString(s.Text),
		Link:        fromPointerToString(s.Link),
	}

	return song, nil
}

func fromPointerToString(ptr *string) string {
	if ptr == nil {
		return ""
	}

	return *ptr
}

func fromPointerToDate(ptr *string) (time.Time, error) {
	if ptr == nil {
		return time.Time{}, nil
	}

	t, err := time.Parse("2006-01-02", *ptr)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

func FromModelToResponse(s *models.Song) *ResponseSong {
	song := &ResponseSong{
		ID:          s.ID.String(),
		Group:       s.GroupName,
		Name:        s.Name,
		ReleaseDate: clearDate(s.ReleaseDate),
		Text:        s.Text,
		Link:        s.Link,
	}

	return song
}

func FromModelsToResponse(modelsSongs []*models.Song) []*ResponseSong {
	songs := make([]*ResponseSong, 0, len(modelsSongs))

	for _, modelsSong := range modelsSongs {
		songs = append(songs, FromModelToResponse(modelsSong))
	}

	return songs
}

func clearDate(t time.Time) string {
	return t.Format("2006-01-02")
}

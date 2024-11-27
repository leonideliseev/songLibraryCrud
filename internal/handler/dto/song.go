package dto

import "github.com/leonideliseev/songLibraryCrud/models"

type RequestCreateSong struct {
	Group       string `json:"group" validate:"required" example:"Imagine Dragons"`
	Name        string `json:"song" validate:"required" example:"Thunder"`
}

type SongDetail struct {
	ReleaseDate string `json:"release_date" validate:"required"`
	Text 		string `json:"text" validate:"required"`
	Link 		string `json:"link" validate:"required"`
}

type ResponseCreateSong struct {
	Song *models.Song `json:"song"`
}

type UpdateSong struct {
	Group       *string `json:"group,omitempty" validate:"omitempty,required" example:"Imagine Dragons"`
	Name        *string `json:"name,omitempty" validate:"omitempty,required" example:"Thunder"`
	ReleaseDate *string `json:"release_date,omitempty" validate:"omitempty,required" example:"27.04.2017"`
	Text 		*string `json:"text,omitempty" validate:"omitempty,required" example:"very much cool song text"`
	Link 		*string `json:"link,omitempty" validate:"omitempty,url" example:"https://www.youtube.com/watch?v=fKopy74weus"`
}

type RequestUpdateSong struct {
}

type ResponseUpdateSong struct {
	Song *models.Song
}

type ResponseGetSongs struct {
	Songs []models.Song
}

type ResponseGetSong struct {
	Song *models.Song
}

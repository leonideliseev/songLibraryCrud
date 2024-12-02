package dto

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
	Song *ResponseSong `json:"song"`
}

type RequestUpdateSong struct {
	Group       *string `json:"group,omitempty" validate:"omitempty,required"`
	Name        *string `json:"name,omitempty" validate:"omitempty,required"`
	ReleaseDate *string `json:"release_date,omitempty" validate:"omitempty,required"`
	Text 		*string `json:"text,omitempty" validate:"omitempty,required"`
	Link 		*string `json:"link,omitempty" validate:"omitempty,url"`
}

type ResponseUpdateSong struct {
	Song *ResponseSong `json:"song"`
}

type ResponseGetSongs struct {
	Songs []*ResponseSong `json:"songs"`
}

type ResponseGetSong struct {
	Song *ResponseSong `json:"song"`
}

type ResponseSong struct {
	ID          string `json:"id"`
	Group       string `json:"group"`
	Name        string `json:"song"`
	ReleaseDate string `json:"release_date"`
	Text 		string `json:"text"`
	Link 		string `json:"link"`
}

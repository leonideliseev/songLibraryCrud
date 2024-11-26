package models

type Song struct {
	ID          string
	Group       string		`gorm:"not null"`
	Name        string		`gorm:"not null"`
	ReleaseDate string		`gorm:"not null"`
	Text 		string		`gorm:"not null"`
	Link 		string		`gorm:"not null"`
}

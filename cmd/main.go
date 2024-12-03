package main

import (
	"github.com/leonideliseev/songLibraryCrud/internal/pkg/app"
)

// @title Song Library API
// @version 1.0
// @description API for managing a library of songs, including creating, retrieving, updating, and deleting songs.

// @contact.name Leonid Eliseev
// @contact.url https://t.me/Lenchiiiikkkk
// @contact.email leonid.2004eliseev@mail.ru

// @BasePath /api/v1
func main() {
	ap := app.NewApp()

	ap.Run()
}

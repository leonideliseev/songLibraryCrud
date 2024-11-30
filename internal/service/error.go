package service

import "errors"

var (
	ErrSongNotFound          = errors.New("song not found")
	ErrSongAlreadyExists     = errors.New("song already exist")
	ErrUpdatedSongNotChanged = errors.New("song for update not changed")
)

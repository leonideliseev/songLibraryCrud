// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package queries

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createSong = `-- name: CreateSong :one
INSERT INTO songs (group_name, name, release_date, text, link)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, group_name, name, release_date, text, link
`

type CreateSongParams struct {
	GroupName   string
	Name        string
	ReleaseDate pgtype.Date
	Text        pgtype.Text
	Link        pgtype.Text
}

func (q *Queries) CreateSong(ctx context.Context, arg CreateSongParams) (Song, error) {
	row := q.db.QueryRow(ctx, createSong,
		arg.GroupName,
		arg.Name,
		arg.ReleaseDate,
		arg.Text,
		arg.Link,
	)
	var i Song
	err := row.Scan(
		&i.ID,
		&i.GroupName,
		&i.Name,
		&i.ReleaseDate,
		&i.Text,
		&i.Link,
	)
	return i, err
}

const deleteSong = `-- name: DeleteSong :exec
DELETE FROM songs
WHERE id = $1
`

func (q *Queries) DeleteSong(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteSong, id)
	return err
}

const getSong = `-- name: GetSong :one
SELECT id, group_name, name, release_date, text, link
FROM songs
WHERE id = $1
`

func (q *Queries) GetSong(ctx context.Context, id pgtype.UUID) (Song, error) {
	row := q.db.QueryRow(ctx, getSong, id)
	var i Song
	err := row.Scan(
		&i.ID,
		&i.GroupName,
		&i.Name,
		&i.ReleaseDate,
		&i.Text,
		&i.Link,
	)
	return i, err
}

const getSongs = `-- name: GetSongs :many
SELECT id, group_name, name, release_date, text, link
FROM songs
LIMIT $1 OFFSET $2
`

type GetSongsParams struct {
	Limit  int32
	Offset int32
}

func (q *Queries) GetSongs(ctx context.Context, arg GetSongsParams) ([]Song, error) {
	rows, err := q.db.Query(ctx, getSongs, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Song
	for rows.Next() {
		var i Song
		if err := rows.Scan(
			&i.ID,
			&i.GroupName,
			&i.Name,
			&i.ReleaseDate,
			&i.Text,
			&i.Link,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateSong = `-- name: UpdateSong :one
UPDATE songs
SET group_name = $2, name = $3, release_date = $4, text = $5, link = $6
WHERE id = $1
RETURNING id, group_name, name, release_date, text, link
`

type UpdateSongParams struct {
	ID          pgtype.UUID
	GroupName   string
	Name        string
	ReleaseDate pgtype.Date
	Text        pgtype.Text
	Link        pgtype.Text
}

func (q *Queries) UpdateSong(ctx context.Context, arg UpdateSongParams) (Song, error) {
	row := q.db.QueryRow(ctx, updateSong,
		arg.ID,
		arg.GroupName,
		arg.Name,
		arg.ReleaseDate,
		arg.Text,
		arg.Link,
	)
	var i Song
	err := row.Scan(
		&i.ID,
		&i.GroupName,
		&i.Name,
		&i.ReleaseDate,
		&i.Text,
		&i.Link,
	)
	return i, err
}
-- name: GetSongs :many
SELECT id, group_name, name, release_date, text, link
FROM songs
LIMIT $1 OFFSET $2;

-- name: CreateSong :one
INSERT INTO songs (group_name, name, release_date, text, link)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetSong :one
SELECT id, group_name, name, release_date, text, link
FROM songs
WHERE id = $1;

-- name: UpdateSong :one
UPDATE songs
SET group_name = $2, name = $3, release_date = $4, text = $5, link = $6
WHERE id = $1
RETURNING *;

-- name: DeleteSong :exec
DELETE FROM songs
WHERE id = $1;

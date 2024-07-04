// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: playlists.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createPlaylist = `-- name: CreatePlaylist :one
INSERT INTO playlists (playlist_id, event_uuid, name)
VALUES ($1, $2, $3)
RETURNING event_uuid, playlist_id, name, created_at, updated_at
`

type CreatePlaylistParams struct {
	PlaylistID string      `json:"playlist_id"`
	EventUuid  pgtype.UUID `json:"event_uuid"`
	Name       string      `json:"name"`
}

func (q *Queries) CreatePlaylist(ctx context.Context, arg CreatePlaylistParams) (Playlist, error) {
	row := q.db.QueryRow(ctx, createPlaylist, arg.PlaylistID, arg.EventUuid, arg.Name)
	var i Playlist
	err := row.Scan(
		&i.EventUuid,
		&i.PlaylistID,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deletePlaylist = `-- name: DeletePlaylist :exec
DELETE FROM playlists
WHERE event_uuid = $1 AND playlist_id = $2
`

type DeletePlaylistParams struct {
	EventUuid  pgtype.UUID `json:"event_uuid"`
	PlaylistID string      `json:"playlist_id"`
}

func (q *Queries) DeletePlaylist(ctx context.Context, arg DeletePlaylistParams) error {
	_, err := q.db.Exec(ctx, deletePlaylist, arg.EventUuid, arg.PlaylistID)
	return err
}

const getAllPlaylists = `-- name: GetAllPlaylists :many
SELECT event_uuid, playlist_id, name, created_at, updated_at
FROM playlists
WHERE event_uuid = $1
`

func (q *Queries) GetAllPlaylists(ctx context.Context, eventUuid pgtype.UUID) ([]Playlist, error) {
	rows, err := q.db.Query(ctx, getAllPlaylists, eventUuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Playlist
	for rows.Next() {
		var i Playlist
		if err := rows.Scan(
			&i.EventUuid,
			&i.PlaylistID,
			&i.Name,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const getPlaylist = `-- name: GetPlaylist :one
SELECT event_uuid, playlist_id, name, created_at, updated_at
FROM playlists
WHERE event_uuid = $1 AND playlist_id = $2
`

type GetPlaylistParams struct {
	EventUuid  pgtype.UUID `json:"event_uuid"`
	PlaylistID string      `json:"playlist_id"`
}

func (q *Queries) GetPlaylist(ctx context.Context, arg GetPlaylistParams) (Playlist, error) {
	row := q.db.QueryRow(ctx, getPlaylist, arg.EventUuid, arg.PlaylistID)
	var i Playlist
	err := row.Scan(
		&i.EventUuid,
		&i.PlaylistID,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getPlaylistUUIDByName = `-- name: GetPlaylistUUIDByName :one
SELECT playlist_id 
FROM playlists
WHERE event_uuid = $1 AND name = $2
`

type GetPlaylistUUIDByNameParams struct {
	EventUuid pgtype.UUID `json:"event_uuid"`
	Name      string      `json:"name"`
}

func (q *Queries) GetPlaylistUUIDByName(ctx context.Context, arg GetPlaylistUUIDByNameParams) (string, error) {
	row := q.db.QueryRow(ctx, getPlaylistUUIDByName, arg.EventUuid, arg.Name)
	var playlist_id string
	err := row.Scan(&playlist_id)
	return playlist_id, err
}

const updatePlaylistName = `-- name: UpdatePlaylistName :one
UPDATE playlists
SET name = $1
WHERE event_uuid = $2 AND playlist_id = $3
RETURNING event_uuid, playlist_id, name, created_at, updated_at
`

type UpdatePlaylistNameParams struct {
	Name       string      `json:"name"`
	EventUuid  pgtype.UUID `json:"event_uuid"`
	PlaylistID string      `json:"playlist_id"`
}

func (q *Queries) UpdatePlaylistName(ctx context.Context, arg UpdatePlaylistNameParams) (Playlist, error) {
	row := q.db.QueryRow(ctx, updatePlaylistName, arg.Name, arg.EventUuid, arg.PlaylistID)
	var i Playlist
	err := row.Scan(
		&i.EventUuid,
		&i.PlaylistID,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
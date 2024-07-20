package events

import (
	"spotify-collab/internal/database"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateEventReq struct {
	UserUUID     uuid.UUID `json:"user_uuid"`
	Name         string    `json:"name"`
	PlaylistUUID uuid.UUID `json:"playlist_uuid"`
}

type CreateEventRes struct {
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	EventUUID uuid.UUID          `json:"event_uuid"`
	Name      string             `json:"name"`
	Playlist  database.Playlist  `json:"playlist"`
}

type ListEventsReq struct {
	UserUUID uuid.UUID `json:"user_uuid"`
}

type GetEventReq struct {
	EventUUID uuid.UUID `json:"event_uuid"`
}

type UpdateEventReq struct {
	EventUUID    uuid.UUID `json:"event_uuid"`
	Name         string    `json:"name"`
	PlaylistUUID uuid.UUID `json:"playlist_uuid"`
}

type DeleteEventReq struct {
	EventUUID uuid.UUID `json:"event_uuid"`
}

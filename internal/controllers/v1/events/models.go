package events

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateEventReq struct {
	UserUUID uuid.UUID `json:"user_uuid"`
	Name     string    `json:"name"`
}

type CreateEventRes struct {
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	EventUUID uuid.UUID          `json:"event_uuid"`
	Name      string             `json:"name"`
}

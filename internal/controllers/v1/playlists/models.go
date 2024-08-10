package playlists

import "github.com/google/uuid"

type CreatePlaylistReq struct {
	UserUUID uuid.UUID `json:"user_uuid" binding:"required"`
	Name     string    `json:"name" binding:"required"`
}

type ListPlaylistsReq struct {
	UserUUID uuid.UUID `json:"user_uuid" binding:"required"`
}

type GetPlaylistReq struct {
	PlaylistUUID uuid.UUID `json:"playlist_uuid" binding:"required"`
}

type UpdatePlaylistReq struct {
	PlaylistUUID uuid.UUID `json:"playlist_uuid" binding:"required"`
	Name         string    `json:"name" binding:"required"`
}

type DeletePlaylistReq struct {
	PlaylistUUID uuid.UUID `json:"playlist_uuid" binding:"required"`
}

type UpdateConfigurationReq struct {
	PlaylistUUID    uuid.UUID `json:"playlist_uuid" binding:"required"`
	Explicit        *bool     `json:"explicit"`
	RequireApproval *bool     `json:"require_approval"`
	MaxSong         int32    `json:"max_song"`
}

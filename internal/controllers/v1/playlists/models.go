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

type CreatePlaylistSpotifyReq struct {
	PlaylistName 	string 	  `json:"playlist_name"`
	IsPublic		bool	  `json:"is_public"`
	IsCollaborative	bool	  `json:"is_collaborative"`
	Description		string	  `json:"description"`

	AccessToken		string	  `json:"access_token"`
	// UserID 			string	  `json:"user_id"`
}

type CreatePlaylistSpotifyReqBody struct {

	PlaylistName 	string 	  `json:"name"`
	IsPublic		bool	  `json:"public"`
	IsCollaborative	bool	  `json:"collaborative"`
	Description		string	  `json:"description"`
}
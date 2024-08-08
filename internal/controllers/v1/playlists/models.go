package playlists

import "github.com/google/uuid"

type CreatePlaylistReq struct {
	UserUUID uuid.UUID `json:"user_uuid"`
	Name     string    `json:"name"`
}

type ListPlaylistsReq struct {
	UserUUID uuid.UUID `json:"user_uuid"`
}

type GetPlaylistReq struct {
	PlaylistUUID uuid.UUID `json:"playlist_uuid"`
}

type UpdatePlaylistReq struct {
	PlaylistUUID uuid.UUID `json:"playlist_uuid"`
	Name         string    `json:"name"`
}

type DeletePlaylistReq struct {
	PlaylistUUID uuid.UUID `json:"playlist_uuid"`
}

type UpdateConfigurationReq struct {
	PlaylistUUID    uuid.UUID `json:"playlist_uuid"`
	Explicit        *bool     `json:"explicit"`
	RequireApproval *bool     `json:"require_approval"`
	MaxSong         *int32    `json:"max_songs"`
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
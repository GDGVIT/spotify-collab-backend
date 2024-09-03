package songs

import "github.com/google/uuid"

type AddSongToPlaylistReq struct {
	PlaylistCode string `json:"playlist_code"`
	SongURI      string `json:"song_uri"`
}

type BlacklistSongReq struct {
	SongURI      string    `json:"song_uri"`
	PlaylistUUID uuid.UUID `json:"playlist_uuid"`
}

type GetAllSongsReq struct {
	PlaylistUUID uuid.UUID `json:"playlist_uuid"`
}

type KaranAddSongToPlaylistReq struct {
	SongURIList []string `json:"uris"`

	// temporary pass playlist id and access token
	PlaylistID  string `json:"playlist_id"`
	AccessToken string `json:"access_token"`
}

type RequestBody struct {
	Uris []string `json:"uris"`
}

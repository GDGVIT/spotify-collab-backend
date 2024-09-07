package songs

import "github.com/google/uuid"

type AddSongToDBReq struct {
	PlaylistCode string `json:"playlist_code"`
	SongURI      string `json:"song_uri"`
}
type AddSongToPlaylistReq struct {
	PlaylistUUID uuid.UUID `json:"playlist_uuid"`
	SongURI      string    `json:"song_uri"`
	Option       string    `uri:"option"`
}

type BlacklistSongReq struct {
	SongURI      string    `json:"song_uri"`
	PlaylistUUID uuid.UUID `json:"playlist_uuid"`
}

type GetAllSongsReq struct {
	PlaylistUUID uuid.UUID `json:"playlist_uuid"`
}

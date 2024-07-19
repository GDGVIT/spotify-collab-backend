package songs

import "github.com/google/uuid"

type AddSongToEventReq struct {
	EventCode string `json:"event_code"`
	SongURI   string `json:"song_uri"`
}

type BlacklistSongReq struct {
	SongURI      string    `json:"song_uri"`
	PlaylistUUID uuid.UUID `json:"playlist_uuid"`
}

type GetAllSongsReq struct {
	PlaylistUUID uuid.UUID `json:"playlist_uuid"`
}

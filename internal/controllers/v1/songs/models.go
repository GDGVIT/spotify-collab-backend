package songs

type AddSongToEventReq struct {
	EventCode string `json:"event_code"`
	URI       string `json:"uri"`
}

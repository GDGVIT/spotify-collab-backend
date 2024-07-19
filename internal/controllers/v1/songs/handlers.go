package songs

import (
	"database/sql"
	"errors"
	"net/http"
	"spotify-collab/internal/database"
	"spotify-collab/internal/merrors"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SongHandler struct {
	db *pgxpool.Pool
}

func Handler(db *pgxpool.Pool) *SongHandler {
	return &SongHandler{
		db: db,
	}
}

func (s *SongHandler) AddSongToEvent(c *gin.Context) {
	req, err := validateAddSongToEventReq(c)
	if err != nil {
		merrors.Validation(c, err.Error())
	}

	q := database.New(s.db)
	event, err := q.GetEventUUIDByCode(c, req.EventCode)
	if errors.Is(sql.ErrNoRows, err) {
		merrors.NotFound(c, "Event not found")
	} else if err != nil {
		merrors.InternalServer(c, err.Error())
	}

	playlist, err := q.GetPlaylistUUIDByEventUUID(c, event)
	if errors.Is(sql.ErrNoRows, err) {
		merrors.NotFound(c, "no playlist found")
	} else if err != nil {
		merrors.InternalServer(c, err.Error())
	}

	// TODO: Check if valid song, passes config -> not greater than count, not blacklisted, other configs
	_, err = q.AddSong(c, database.AddSongParams{
		SongUri:      req.URI,
		PlaylistUuid: playlist,
	})
	if err != nil {
		merrors.InternalServer(c, err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "song added",
	})
}

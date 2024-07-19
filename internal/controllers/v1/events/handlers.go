package events

import (
	"database/sql"
	"errors"
	"math/rand"
	"net/http"
	"spotify-collab/internal/database"
	"spotify-collab/internal/merrors"
	"spotify-collab/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EventHandler struct {
	db *pgxpool.Pool
}

func Handler(db *pgxpool.Pool) *EventHandler {
	return &EventHandler{
		db: db,
	}
}

func (e *EventHandler) CreateEvent(c *gin.Context) {
	req, err := validateCreateEventReq(c)
	if err != nil {
		merrors.Validation(c, err.Error())
		return
	}

	tx, err := e.db.Begin(c)
	if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}
	defer tx.Rollback(c)
	qtx := database.New(e.db).WithTx(tx)

	// TODO: Recheck if duplicate
	eventCode := GenerateEventCode(6)

	event, err := qtx.CreateEvent(c, database.CreateEventParams{
		UserUuid:  req.UserUUID,
		Name:      req.Name,
		EventCode: eventCode,
	})
	if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}

	playlist, err := qtx.GetPlaylist(c, req.PlaylistUUID)
	if errors.Is(sql.ErrNoRows, err) {
		merrors.NotFound(c, "Playlist not found")
		return
	} else if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}

	err = qtx.SetPlaylistForEvent(c, database.SetPlaylistForEventParams{
		EventUuid:    event.EventUuid,
		PlaylistUuid: req.PlaylistUUID,
	})
	if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.BaseResponse{
		Success:    true,
		StatusCode: http.StatusOK,
		Message:    "event created successfully",
		Data: CreateEventRes{
			CreatedAt: event.CreatedAt,
			EventUUID: event.EventUuid,
			Name:      event.Name,
			Playlist:  playlist,
		},
	})
}

func (e *EventHandler) ListEvents(c *gin.Context) {
	req, err := validateListEventsReq(c)
	if err != nil {
		merrors.Validation(c, err.Error())
		return
	}

	q := database.New(e.db)
	events, err := q.GetAllEvents(c, req.UserUUID)
	if errors.Is(sql.ErrNoRows, err) {
		merrors.NotFound(c, "No Events exist!")
	} else if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.BaseResponse{
		Success:    true,
		Message:    "Events successfully retrieved",
		Data:       events,
		StatusCode: http.StatusOK,
	})
}

func (e *EventHandler) GetEvent(c *gin.Context) {
	req, err := validateGetEventReq(c)
	if err != nil {
		merrors.Validation(c, err.Error())
		return
	}

	q := database.New(e.db)
	event, err := q.GetEvent(c, req.EventUUID)
	if errors.Is(sql.ErrNoRows, err) {
		merrors.NotFound(c, "Event not found")
	} else if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.BaseResponse{
		Success:    true,
		Message:    "Event successfully retrieved",
		Data:       event,
		StatusCode: http.StatusOK,
	})
}

func (e *EventHandler) UpdateEvent(c *gin.Context) {
	req, err := validateUpdateEventReq(c)
	if err != nil {
		merrors.Validation(c, err.Error())
		return
	}

	q := database.New(e.db)
	event, err := q.UpdateEventName(c, database.UpdateEventNameParams{
		Name:      req.Name,
		EventUuid: req.EventUUID,
	})
	if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.BaseResponse{
		Success:    true,
		Message:    "Event successfully updated",
		Data:       event,
		StatusCode: http.StatusOK,
	})
}

func (e *EventHandler) DeleteEvent(c *gin.Context) {
	req, err := validateDeleteEventReq(c)
	if err != nil {
		merrors.Validation(c, err.Error())
		return
	}

	q := database.New(e.db)
	rows, err := q.DeleteEvent(c, req.EventUUID)
	if rows == 0 {
		merrors.NotFound(c, "Event not found")
	} else if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, utils.BaseResponse{
		Success:    true,
		Message:    "Event successfully deleted",
		StatusCode: http.StatusOK,
	})
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func GenerateEventCode(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

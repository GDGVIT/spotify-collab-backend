package events

import (
	"math/rand"
	"net/http"
	"spotify-collab/internal/database"
	"spotify-collab/internal/merrors"

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

	q := database.New(e.db)

	eventCode := GenerateEventCode(6)

	event, err := q.CreateEvent(c, database.CreateEventParams{
		UserUuid:  req.UserUUID,
		Name:      req.Name,
		EventCode: eventCode,
	})
	if err != nil {
		merrors.InternalServer(c, err.Error())
	}

	c.JSON(http.StatusOK, CreateEventRes{
		CreatedAt: event.CreatedAt,
		EventUUID: event.EventUuid,
		Name:      event.Name,
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

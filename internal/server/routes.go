package server

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"spotify-collab/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.GET("/", s.HelloWorldHandler)

	r.GET("/health", s.healthHandler)

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}

func (s *Server) AddSongToEvent(c *gin.Context) {
	var input struct {
		EventCode string `json:"event_code"`
		URI       string `json:"uri"`
	}

	if err := c.ShouldBindJSON(input); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"msg": "invalid format",
		})
	}

	q := database.New(s.db.Db)
	event, err := q.GetEventUUIDByCode(c, input.EventCode)
	if errors.Is(sql.ErrNoRows, err) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "event not found",
		})
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err,
		})
	}

	playlist, err := q.GetPlaylistUUIDByEventUUID(c, event)
	if errors.Is(sql.ErrNoRows, err) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "no playlist found",
		})
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err,
		})
	}

	// TODO: Check if valid song, passes config -> not greater than count, not blacklisted, other configs
	_, err = q.AddSong(c, database.AddSongParams{
		SongUri:    input.URI,
		PlaylistID: playlist,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "song added",
	})
}

func (s *Server) CreateEvent(c *gin.Context) {
	var input struct {
		User_uuid uuid.UUID `json:"user_uuid"`
		Name      string    `json:"name"`
	}

	if err := c.ShouldBindJSON(input); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"msg": "invalid format",
		})
	}

	q := database.New(s.db.Db)

	eventCode := GenerateEventCode(6)

	event, err := q.CreateEvent(c, database.CreateEventParams{
		UserUuid:  input.User_uuid,
		Name:      input.Name,
		EventCode: eventCode,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": fmt.Errorf("server error: %w", err).Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"created_at": event.CreatedAt,
		"event_uuid": event.EventUuid,
		"name":       event.Name,
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

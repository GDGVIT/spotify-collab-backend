package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"

	"spotify-collab/internal/controllers/v1/events"
	"spotify-collab/internal/controllers/v1/playlists"
	"spotify-collab/internal/controllers/v1/songs"
	"spotify-collab/internal/database"
)

type Server struct {
	port int

	db              *pgxpool.Pool
	eventHandler    *events.EventHandler
	playlistHandler *playlists.PlaylistHandler
	songHandler     *songs.SongHandler
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	db := database.NewService()
	NewServer := &Server{
		port: port,

		db:              db,
		eventHandler:    events.Handler(db),
		playlistHandler: playlists.Handler(db),
		songHandler:     songs.Handler(db),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf("localhost:%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

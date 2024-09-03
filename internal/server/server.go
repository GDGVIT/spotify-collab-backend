package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	spotifyauth "github.com/zmb3/spotify/v2/auth"

	"spotify-collab/internal/controllers/v1/auth"
	"spotify-collab/internal/controllers/v1/events"
	"spotify-collab/internal/controllers/v1/playlists"
	"spotify-collab/internal/controllers/v1/songs"
	"spotify-collab/internal/database"
)

type Server struct {
	port int

	db          *pgxpool.Pool
	spotifyauth *spotifyauth.Authenticator

	eventHandler    *events.EventHandler
	playlistHandler *playlists.PlaylistHandler
	songHandler     *songs.SongHandler
	authHandler     *auth.AuthHandler
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	db := database.NewService()

	redirectURI := "http://localhost:8080/auth/spotify/callback"
	scopes := []string{
		spotifyauth.ScopeUserModifyPlaybackState,
		spotifyauth.ScopePlaylistModifyPrivate,
		spotifyauth.ScopePlaylistModifyPublic,
		spotifyauth.ScopeUserReadPrivate,
		spotifyauth.ScopeUserReadCurrentlyPlaying,
		spotifyauth.ScopeUserReadPlaybackState,
	}
	spotifyAuth := spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(scopes...))

	NewServer := &Server{
		port: port,

		db:          db,
		spotifyauth: spotifyAuth,

		eventHandler:    events.Handler(db),
		playlistHandler: playlists.Handler(db, spotifyAuth),
		songHandler:     songs.Handler(db),
		authHandler:     auth.Handler(db, spotifyAuth),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

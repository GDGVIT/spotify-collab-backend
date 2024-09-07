package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()
	r.Use(CORSMiddleware())

	r.GET("/", s.HelloWorldHandler)

	r.GET("/health", s.healthHandler)

	v1 := r.Group("/v1")

	v1.POST("/playlists", s.authenticate(), s.playlistHandler.CreatePlaylist)
	v1.GET("/playlists", s.authenticate(), s.playlistHandler.ListPlaylists)
	v1.GET("/playlists/:id", s.authenticate(), s.playlistHandler.GetPlaylist)
	v1.POST("/playlists/:id", s.authenticate(), s.playlistHandler.UpdatePlaylist)
	v1.DELETE("/playlists/:id", s.authenticate(), s.playlistHandler.DeletePlaylist)

	v1.POST("/songs/add", s.songHandler.AddSongToDB)                             // Called by the participants
	v1.POST("/songs/:option", s.authenticate(), s.songHandler.AddSongToPlaylist) // Either accept or reject, called by the admin to add it to playlist or reject
	v1.GET("/songs", s.authenticate(), s.songHandler.GetAllSongs)

	v1.POST("/songs/blacklist", s.songHandler.BlacklistSong)
	v1.GET("/songs/blacklist", s.songHandler.GetBlacklistedSongs)
	v1.DELETE("/songs/blacklist", s.songHandler.DeleteBlacklistSong)

	auth := v1.Group("/auth")
	auth.GET("spotify/login", s.authHandler.SpotifyLogin)
	auth.GET("spotify/callback", s.authHandler.SpotifyCallback)

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	err := s.db.Ping(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf(fmt.Sprintf("db down: %v", err)) // Log the error and terminate the program
		c.JSON(http.StatusInternalServerError, stats)
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"
	c.JSON(http.StatusOK, stats)
}

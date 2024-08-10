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

	r.GET("/", s.HelloWorldHandler)

	r.GET("/health", s.healthHandler)

	r.POST("/events/new", s.eventHandler.CreateEvent)
	r.GET("/events/list", s.eventHandler.ListEvents)
	r.GET("/events/one", s.eventHandler.GetEvent)
	r.POST("/events/one", s.eventHandler.UpdateEvent)
	r.DELETE("/events/", s.eventHandler.DeleteEvent)

	r.POST("/playlists", s.playlistHandler.CreatePlaylist)
	r.GET("/playlists/list", s.playlistHandler.ListPlaylists)
	r.GET("/playlists/one", s.playlistHandler.GetPlaylist)
	r.POST("/playlists/one", s.playlistHandler.UpdatePlaylist)
	r.DELETE("/playlists/", s.playlistHandler.DeletePlaylist)
	r.PATCH("/playlists/config", s.playlistHandler.UpdateConfiguration)

	r.POST("/songs/new", s.songHandler.AddSongToEvent)
	r.POST("/songs/blacklist", s.songHandler.BlacklistSong)
	r.GET("/songs/all", s.songHandler.GetAllSongs)
	r.GET("/songs/blacklist", s.songHandler.GetBlacklistedSongs)
	r.DELETE("/songs/blacklist", s.songHandler.DeleteBlacklistSong)

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

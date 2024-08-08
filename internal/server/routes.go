package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.GET("/", s.HelloWorldHandler)

	r.GET("/health", s.healthHandler)

	r.POST("/events/new", s.eventHandler.CreateEvent)

	r.POST("/playlists", s.playlistHandler.CreatePlaylist)
	r.GET("/playlists", s.playlistHandler.ListPlaylists)
	r.GET("/playlists/:id", s.playlistHandler.GetPlaylist)
	r.PATCH("/playlists/:id", s.playlistHandler.UpdatePlaylist)
	r.DELETE("/playlists/:id", s.playlistHandler.DeletePlaylist)

	r.POST("/songs/new", s.songHandler.AddSongToEvent)
	
	// Route needs to be changed
	r.POST("/songs/add", s.songHandler.AddSongToPlaylist)

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	// c.JSON(http.StatusOK, s.db.Health())
	c.JSON(http.StatusOK, gin.H{
		"testing": "ready!",
	})
}

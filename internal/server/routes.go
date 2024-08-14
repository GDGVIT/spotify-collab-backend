package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"spotify-collab/internal/merrors"
	"spotify-collab/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.GET("/", s.HelloWorldHandler)

	r.GET("/health", s.healthHandler)
	r.GET("/schema", s.schemaHandler)
	r.GET("/test", s.testHandler)

	r.POST("/events/new", s.eventHandler.CreateEvent)
	r.GET("/events/list", s.eventHandler.ListEvents)
	r.GET("/events/one", s.eventHandler.GetEvent)
	r.POST("/events/one", s.eventHandler.UpdateEvent)
	r.DELETE("/events/", s.eventHandler.DeleteEvent)

	r.POST("/playlists", s.playlistHandler.CreatePlaylist)
	r.GET("/playlists", s.playlistHandler.ListPlaylists)
	r.GET("/playlists/:id", s.playlistHandler.GetPlaylist)
	r.POST("/playlists/:id", s.playlistHandler.UpdatePlaylist)
	r.DELETE("/playlists/:id", s.playlistHandler.DeletePlaylist)
	r.PATCH("/playlists/config", s.playlistHandler.UpdateConfiguration)

	r.POST("/songs/new", s.songHandler.AddSongToEvent)
	r.POST("/songs/blacklist", s.songHandler.BlacklistSong)
	r.GET("/songs/all", s.songHandler.GetAllSongs)
	r.GET("/songs/blacklist", s.songHandler.GetBlacklistedSongs)
	r.DELETE("/songs/blacklist", s.songHandler.DeleteBlacklistSong)

	// Route needs to be changed
	r.POST("/songs/add", s.songHandler.AddSongToPlaylist)
	r.POST("/playlists/add", s.playlistHandler.CreatePlaylistSpotify)

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

func (s *Server) schemaHandler(c *gin.Context) {
	data := gin.H{}

	rows, err := s.db.Query(context.Background(), "SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname != 'pg_catalog' AND schemaname != 'information_schema'")
	if err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}

	defer rows.Close()
	var tables []string
	for rows.Next() {
		var i string
		if err := rows.Scan(
			&i,
		); err != nil {
			merrors.InternalServer(c, err.Error())
			return
		}
		tables = append(tables, i)
	}
	if err := rows.Err(); err != nil {
		merrors.InternalServer(c, err.Error())
		return
	}
	data["tables"] = tables

	rows, err = s.db.Query(context.Background(), "select column_name from INFORMATION_SCHEMA.COLUMNS where table_name = 'users';")
	if err != nil {
		data["error"] = err.Error()
		c.JSON(http.StatusInternalServerError, utils.BaseResponse{
			Success:    true,
			Data:       data,
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	defer rows.Close()
	var details []string
	for rows.Next() {
		var i string
		if err := rows.Scan(
			&i,
		); err != nil {
			data["error"] = err.Error()
			c.JSON(http.StatusInternalServerError, utils.BaseResponse{
				Success:    true,
				Data:       data,
				StatusCode: http.StatusInternalServerError,
			})
			return
		}
		details = append(details, i)
	}
	if err := rows.Err(); err != nil {
		data["error"] = err.Error()
		c.JSON(http.StatusInternalServerError, utils.BaseResponse{
			Success:    true,
			Data:       data,
			StatusCode: http.StatusInternalServerError,
		})
		return
	}
	data["details"] = details

	rows, err = s.db.Query(context.Background(), "select name from users;")
	if err != nil {
		data["error"] = err.Error()
		c.JSON(http.StatusInternalServerError, utils.BaseResponse{
			Success:    true,
			Data:       data,
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	defer rows.Close()
	var users []string
	for rows.Next() {
		var i string
		if err := rows.Scan(
			&i,
		); err != nil {
			data["error"] = err.Error()
			c.JSON(http.StatusInternalServerError, utils.BaseResponse{
				Success:    true,
				Data:       data,
				StatusCode: http.StatusInternalServerError,
			})
			return
		}
		users = append(users, i)
	}
	if err := rows.Err(); err != nil {
		data["error"] = err.Error()
		c.JSON(http.StatusInternalServerError, utils.BaseResponse{
			Success:    true,
			Data:       data,
			StatusCode: http.StatusInternalServerError,
		})
		return
	}
	data["users"] = users

	c.JSON(http.StatusOK, utils.BaseResponse{
		Success:    true,
		Data:       data,
		StatusCode: http.StatusOK,
	})

}

func (s *Server) testHandler(c *gin.Context) {
	var dbName, schemaName string
	_ = s.db.QueryRow(context.Background(), "SELECT current_database(), current_schema()").Scan(&dbName, &schemaName)
	log.Printf("Current DB: %s, Current Schema: %s\n", dbName, schemaName)
}

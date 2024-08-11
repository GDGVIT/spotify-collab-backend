package playlists

import (
	"github.com/gin-gonic/gin"
)

func validateCreatePlaylist(c *gin.Context) (CreatePlaylistReq, error) {
	var req CreatePlaylistReq
	err := c.ShouldBindJSON(&req)
	return req, err
}

func validateListPlaylistsReq(c *gin.Context) (ListPlaylistsReq, error) {
	var req ListPlaylistsReq
	err := c.ShouldBindJSON(&req)
	return req, err
}
func validateGetPlaylistReq(c *gin.Context) (GetPlaylistReq, error) {
	var req GetPlaylistReq
	err := c.ShouldBindUri(&req)
	return req, err
}

func validateUpdatePlaylistReq(c *gin.Context) (UpdatePlaylistReq, error) {
	var req UpdatePlaylistReq
	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		return req, err
	}
	err = c.ShouldBindUri(&req)
	return req, err
}
func validateDeletePlaylistReq(c *gin.Context) (DeletePlaylistReq, error) {
	var req DeletePlaylistReq
	err := c.ShouldBindUri(&req)
	return req, err
}

func validateUpdateConfigurationReq(c *gin.Context) (UpdateConfigurationReq, error) {
	var req UpdateConfigurationReq
	err := c.ShouldBindJSON(&req)
	return req, err
}

func validateCreatePlaylistSpotifyReq(c *gin.Context) (CreatePlaylistSpotifyReq, error) {
	var req CreatePlaylistSpotifyReq
	err := c.ShouldBindJSON(&req)
	return req, err
}

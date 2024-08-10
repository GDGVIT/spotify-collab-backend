package songs

import (
	"github.com/gin-gonic/gin"
)

func validateAddSongToEventReq(c *gin.Context) (AddSongToEventReq, error) {
	var req AddSongToEventReq
	err := c.ShouldBindJSON(req)
	return req, err
}

func validateBlacklistSongReq(c *gin.Context) (BlacklistSongReq, error) {
	var req BlacklistSongReq
	err := c.ShouldBindJSON(req)
	return req, err
}

func validateGetAllSongsReq(c *gin.Context) (GetAllSongsReq, error) {
	var req GetAllSongsReq
	err := c.ShouldBindJSON(req)
	return req, err
}

func validateAddSongToPlaylist(c *gin.Context) (*AddSongToPlaylistReq, error) {
    var req AddSongToPlaylistReq
    if err := c.ShouldBindJSON(&req); err != nil {
        return nil, err
    }
    return &req, nil
}


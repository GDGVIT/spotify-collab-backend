package songs

import (
	"github.com/gin-gonic/gin"
)

func validateAddSongToEventReq(c *gin.Context) (AddSongToEventReq, error) {
	var req AddSongToEventReq
	err := c.ShouldBindJSON(req)
	return req, err
}

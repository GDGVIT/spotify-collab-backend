package events

import (
	"github.com/gin-gonic/gin"
)

func validateCreateEventReq(c *gin.Context) (CreateEventReq, error) {
	var req CreateEventReq
	err := c.ShouldBindJSON(req)
	return req, err
}

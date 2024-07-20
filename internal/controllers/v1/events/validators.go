package events

import (
	"github.com/gin-gonic/gin"
)

func validateCreateEventReq(c *gin.Context) (CreateEventReq, error) {
	var req CreateEventReq
	err := c.ShouldBindJSON(req)
	return req, err
}

func validateListEventsReq(c *gin.Context) (ListEventsReq, error) {
	var req ListEventsReq
	err := c.ShouldBindJSON(req)
	return req, err
}

func validateGetEventReq(c *gin.Context) (GetEventReq, error) {
	var req GetEventReq
	err := c.ShouldBindJSON(req)
	return req, err
}

func validateUpdateEventReq(c *gin.Context) (UpdateEventReq, error) {
	var req UpdateEventReq
	err := c.ShouldBindJSON(req)
	return req, err
}

func validateDeleteEventReq(c *gin.Context) (DeleteEventReq, error) {
	var req DeleteEventReq
	err := c.ShouldBindJSON(req)
	return req, err
}

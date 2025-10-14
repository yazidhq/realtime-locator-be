package websocket

import (
	"TeamTrackerBE/internal/utils/responses"

	"github.com/gin-gonic/gin"
)

func LiveTrack(c *gin.Context) {
	groupID := c.Query("group_id")
	userID := c.Query("user_id")

	if groupID == "" || userID == "" {
		responses.NewBadRequestError("group id and user id is required")
	}

	ServeWs(c.Writer, c.Request, groupID, userID)
}
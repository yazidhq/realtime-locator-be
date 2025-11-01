package websocket

import (
	"TeamTrackerBE/internal/domain/repository"
	"TeamTrackerBE/internal/utils/responses"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type LiveTrackHandler struct {
	repoUser repository.UserRepository
}

func NewLiveTrackHandler(ru *repository.UserRepository) *LiveTrackHandler {
	return &LiveTrackHandler{ repoUser: *ru }
}

func (u *LiveTrackHandler) LiveTrack(c *gin.Context) {
	userID := c.Query("user_id")

	userIDtoken, _ := c.Get("userID")
	if userIDtoken != userID {
		responses.Error(c, http.StatusBadRequest, "user id and user id from token not match")
		return
	}

	if userID == "" {
		responses.Error(c, http.StatusBadRequest, "user id is required")
		return
	}
	
	userIDParse, errUserID := uuid.Parse(userID);
	if errUserID != nil {
		responses.Error(c, http.StatusBadRequest, "invalid user id, must be uuid type")
		return
	}

	user, err := u.repoUser.FindById(userIDParse)
	if err != nil {
		responses.Error(c, http.StatusBadRequest, "user id not found in user")
		return
	}

	if websocket.IsWebSocketUpgrade(c.Request) {
		ServeWs(c.Writer, c.Request, user)
		return
	}

	responses.Success(c, "All good mate, NEXT CONNECT!", nil)
}
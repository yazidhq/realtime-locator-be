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

	userIDtoken, exists := c.Get("userID")
	if !exists {
		responses.Error(c, http.StatusUnauthorized, "token user not found")
		return
	}

	if userIDParse != userIDtoken {
		responses.Error(c, http.StatusBadRequest, "user id and user id from token not match")
		return
	}

	if websocket.IsWebSocketUpgrade(c.Request) {
		ServeWs(c.Writer, c.Request, user)
		return
	}

	responses.Success(c, "All good mate, NEXT CONNECT!", nil)
}

func (h *LiveTrackHandler) GetUserOnlineStatus(c *gin.Context) {
    idStr := c.Param("id")

    uid, err := uuid.Parse(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
        return
    }

	online := GetHub().IsOnline(uid)
	c.JSON(http.StatusOK, gin.H{
		"user_id": uid,
		"online":  online,
	})
}

func (h *Hub) IsOnline(userID uuid.UUID) bool {
    return h.Online[userID]
}
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
	repoGroup repository.GroupRepository
	repoUser repository.UserRepository
	repoGroupParticipant repository.GroupParticipantRepository
}

func NewLiveTrackHandler(
	rg *repository.GroupRepository,
	ru *repository.UserRepository,
	rgp *repository.GroupParticipantRepository,
	) *LiveTrackHandler {
	return &LiveTrackHandler{
		repoGroup: *rg,
		repoUser: *ru,
		repoGroupParticipant: *rgp,
	}
}

func (u *LiveTrackHandler) LiveTrack(c *gin.Context) {
	groupID := c.Query("group_id")
	userID := c.Query("user_id")

	if groupID == "" || userID == "" {
		responses.Error(c, http.StatusBadRequest, "group id and user id is required")
		return
	}

	groupIDParse, errGroupID := uuid.Parse(groupID);
	if errGroupID != nil {
		responses.Error(c, http.StatusBadRequest, "invalid group id, must be uuid type")
		return
	}

	if _, err := u.repoGroup.FindById(groupIDParse); err != nil {
		responses.Error(c, http.StatusBadRequest, "group id not found in group")
		return
	}
	
	userIDParse, errUserID := uuid.Parse(userID);
	if errUserID != nil {
		responses.Error(c, http.StatusBadRequest, "invalid user id, must be uuid type")
		return
	}

	if _, err := u.repoUser.FindById(userIDParse); err != nil {
		responses.Error(c, http.StatusBadRequest, "user id not found in user")
		return
	}

	if _, err := u.repoGroupParticipant.FindByGroupIDUserID(groupIDParse, userIDParse); err != nil {
		responses.Error(c, http.StatusBadRequest, "user is not group participant")
		return
	}

	if websocket.IsWebSocketUpgrade(c.Request) {
		ServeWs(c.Writer, c.Request, groupID, userID)
		return
	}

	responses.Success(c, "All good mate, NEXT CONNECT!", nil)
}
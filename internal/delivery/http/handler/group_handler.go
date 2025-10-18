package handler

import (
	"TeamTrackerBE/internal/delivery/http/dto"
	"TeamTrackerBE/internal/domain/model"
	"TeamTrackerBE/internal/utils"
	"TeamTrackerBE/internal/utils/responses"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GroupInterface interface {
	Create(group *dto.GroupCreateRequest) (*model.Group, error)
	Update(groupID uuid.UUID, req *dto.GroupUpdateRequest) (*model.Group, error)
	Delete(groupID uuid.UUID) (*model.Group, error)
	FindAll(page, limit int, filters []utils.FilterOptions) ([]model.Group, int, error)
	FindById(groupID uuid.UUID) (*model.Group, error)
	Truncate() (error)
}

type GroupHandler struct {
	uc GroupInterface
}

func NewGroupHandler(uc GroupInterface) *GroupHandler {
	return &GroupHandler{uc: uc}
}

func (h GroupHandler) Create(c *gin.Context) {
	var req dto.GroupCreateRequest
	if errReq := c.ShouldBindJSON(&req); errReq != nil {
		responses.Error(c, http.StatusBadRequest, "Invalid request body: " + errReq.Error())
		return
	}

	created, err := h.uc.Create(&req)
	if err != nil {
		if ce, ok := err.(responses.CodedError); ok {
			responses.Error(c, ce.StatusCode(), ce.Error())
		} else {
			responses.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	responses.Created(c, "Created successfully", created)
}

func (h *GroupHandler) Update(c *gin.Context) {
    id := c.Param("id")

	groupID, err := uuid.Parse(id)
    if err != nil {
        responses.Error(c, http.StatusBadRequest, "Invalid UUID")
        return
    }

    var req dto.GroupUpdateRequest
    if errReq := c.ShouldBindJSON(&req); errReq != nil {
        responses.Error(c, http.StatusBadRequest, "Invalid request body: " + errReq.Error())
        return
    }

    updated, errUpdate := h.uc.Update(groupID, &req)
    if errUpdate != nil {
        if ce, ok := errUpdate.(responses.CodedError); ok {
			responses.Error(c, ce.StatusCode(), ce.Error())
		} else {
			responses.Error(c, http.StatusInternalServerError, errUpdate.Error())
		}
		return
    }

    responses.Success(c, "Updated successfully", updated)
}

func (h GroupHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	groupID, err := uuid.Parse(id)
    if err != nil {
        responses.Error(c, http.StatusBadRequest, "Invalid UUID")
        return
    }

	_, errDelete := h.uc.Delete(groupID)
	if errDelete != nil {
		if ce, ok := errDelete.(responses.CodedError); ok {
			responses.Error(c, ce.StatusCode(), ce.Error())
		} else {
			responses.Error(c, http.StatusInternalServerError, errDelete.Error())
		}
		return
	}
	
	responses.Success(c, "Deleted successfully", nil)
}

func (h GroupHandler) FindAll(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

    if page <= 0 {
        page = 1
    }
    if limit <= 0 {
        limit = 10
    }

    allowedOps := []string{"=", "like", ">", "<"}
    filters := utils.BuildDynamicFilters(c.Request.URL.Query(), allowedOps)

    groups, total, err := h.uc.FindAll(page, limit, filters)
    if err != nil {
        responses.Error(c, http.StatusInternalServerError, err.Error())
        return
    }

    var response []dto.GroupResponse
    for _, u := range groups {
        response = append(response, dto.GroupResponse{
            ID: u.ID.String(),
            Name: u.Name,
            OwnerID: u.OwnerID,
            RadiusArea: u.RadiusArea,
        })
    }

    responses.SuccessPaginated(c, "Get data successfully", response, page, limit, total)
}

func (h GroupHandler) FindById(c *gin.Context) {
	id := c.Param("id")
	groupID, err := uuid.Parse(id)
    if err != nil {
        responses.Error(c, http.StatusBadRequest, "Invalid UUID")
        return
    }

	result, errResult := h.uc.FindById(groupID)
	if errResult != nil {
		if ce, ok := errResult.(responses.CodedError); ok {
			responses.Error(c, ce.StatusCode(), ce.Error())
		} else {
			responses.Error(c, http.StatusInternalServerError, errResult.Error())
		}
		return
	}

	response := dto.GroupResponse {
        ID: result.ID.String(),
		Name: result.Name,
        OwnerID: result.OwnerID,
        RadiusArea: result.RadiusArea,
    }

	responses.Success(c, "Get data successfully", response)
}

func (h *GroupHandler) Truncate(c *gin.Context) {
	err := h.uc.Truncate()
	if err != nil {
		if ce, ok := err.(responses.CodedError); ok {
			responses.Error(c, ce.StatusCode(), ce.Error())
		} else {
			responses.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	responses.Success(c, "Truncated successfully", nil)
}

func (h GroupHandler) Invite(c *gin.Context) {
	id := c.Param("id")
	groupID, err := uuid.Parse(id)
	if err != nil {
		responses.Error(c, http.StatusBadRequest, "Invalid UUID")
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists {
		responses.Error(c, http.StatusUnauthorized, "unauthenticated")
		return
	}

	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		responses.Error(c, http.StatusUnauthorized, "invalid user id in context")
		return
	}

	group, err := h.uc.FindById(groupID)
	if err != nil {
		if ce, ok := err.(responses.CodedError); ok {
			responses.Error(c, ce.StatusCode(), ce.Error())
		} else {
			responses.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	if group.OwnerID != userID {
		responses.Error(c, http.StatusForbidden, "only group owner can create invite link")
		return
	}

	token, err := utils.GenerateGroupInviteToken(groupID)
	if err != nil {
		responses.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	inviteToken := token

	responses.Success(c, "Invite token generated", gin.H{"invite_token": inviteToken})
}
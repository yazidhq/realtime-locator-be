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

type GroupParticipantInterface interface {
	Create(groupParticipant *dto.GroupParticipantCreateRequest) (*model.GroupParticipant, error)
	Update(groupParticipantID uuid.UUID, req *dto.GroupParticipantUpdateRequest) (*model.GroupParticipant, error)
	Delete(groupParticipantID uuid.UUID) (*model.GroupParticipant, error)
	FindAll(page, limit int, filters []utils.FilterOptions) ([]model.GroupParticipant, int, error)
	FindById(groupParticipantID uuid.UUID) (*model.GroupParticipant, error)
	Truncate() (error)
}

type GroupParticipantHandler struct {
	uc GroupParticipantInterface
}

func NewGroupParticipantHandler(uc GroupParticipantInterface) *GroupParticipantHandler {
	return &GroupParticipantHandler{uc: uc}
}

func (h GroupParticipantHandler) Create(c *gin.Context) {
	var req dto.GroupParticipantCreateRequest
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

func (h *GroupParticipantHandler) Update(c *gin.Context) {
    id := c.Param("id")

	groupParticipantID, err := uuid.Parse(id)
    if err != nil {
        responses.Error(c, http.StatusBadRequest, "Invalid UUID")
        return
    }

    var req dto.GroupParticipantUpdateRequest
    if errReq := c.ShouldBindJSON(&req); errReq != nil {
        responses.Error(c, http.StatusBadRequest, "Invalid request body: " + errReq.Error())
        return
    }

    updated, errUpdate := h.uc.Update(groupParticipantID, &req)
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

func (h GroupParticipantHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	groupParticipantID, err := uuid.Parse(id)
    if err != nil {
        responses.Error(c, http.StatusBadRequest, "Invalid UUID")
        return
    }

	_, errDelete := h.uc.Delete(groupParticipantID)
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

func (h GroupParticipantHandler) FindAll(c *gin.Context) {
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

    groupParticipants, total, err := h.uc.FindAll(page, limit, filters)
    if err != nil {
        responses.Error(c, http.StatusInternalServerError, err.Error())
        return
    }

    var response []dto.GroupParticipantResponse
    for _, u := range groupParticipants {
        response = append(response, dto.GroupParticipantResponse{
            ID: u.ID.String(),
            GroupID: u.GroupID,
            UserID: u.UserID,
        })
    }

    responses.SuccessPaginated(c, "Get data successfully", response, page, limit, total)
}

func (h GroupParticipantHandler) FindById(c *gin.Context) {
	id := c.Param("id")
	groupParticipantID, err := uuid.Parse(id)
    if err != nil {
        responses.Error(c, http.StatusBadRequest, "Invalid UUID")
        return
    }

	result, errResult := h.uc.FindById(groupParticipantID)
	if errResult != nil {
		if ce, ok := errResult.(responses.CodedError); ok {
			responses.Error(c, ce.StatusCode(), ce.Error())
		} else {
			responses.Error(c, http.StatusInternalServerError, errResult.Error())
		}
		return
	}

	response := dto.GroupParticipantResponse {
        ID: result.ID.String(),
        GroupID: result.GroupID,
		UserID: result.UserID,
    }

	responses.Success(c, "Get data successfully", response)
}

func (h *GroupParticipantHandler) Truncate(c *gin.Context) {
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
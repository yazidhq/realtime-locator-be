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

type ContactInterface interface {
	Create(contact *dto.ContactCreateRequest) (*model.Contact, error)
	Update(contactID uuid.UUID, req *dto.ContactUpdateRequest) (*model.Contact, error)
	Delete(contactID uuid.UUID) (*model.Contact, error)
	FindAll(page, limit int, filters []utils.FilterOptions) ([]model.Contact, int, error)
	FindById(contactID uuid.UUID) (*model.Contact, error)
	Truncate() (error)
}

type ContactHandler struct {
	uc ContactInterface
}

func NewContactHandler(uc ContactInterface) *ContactHandler {
	return &ContactHandler{uc: uc}
}

func (h ContactHandler) Create(c *gin.Context) {
	var req dto.ContactCreateRequest
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

func (h *ContactHandler) Update(c *gin.Context) {
    id := c.Param("id")

	contactID, err := uuid.Parse(id)
    if err != nil {
        responses.Error(c, http.StatusBadRequest, "Invalid UUID")
        return
    }

    var req dto.ContactUpdateRequest
    if errReq := c.ShouldBindJSON(&req); errReq != nil {
        responses.Error(c, http.StatusBadRequest, "Invalid request body: " + errReq.Error())
        return
    }

    updated, errUpdate := h.uc.Update(contactID, &req)
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

func (h ContactHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	contactID, err := uuid.Parse(id)
    if err != nil {
        responses.Error(c, http.StatusBadRequest, "Invalid UUID")
        return
    }

	_, errDelete := h.uc.Delete(contactID)
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

func (h ContactHandler) FindAll(c *gin.Context) {
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

    contacts, total, err := h.uc.FindAll(page, limit, filters)
    if err != nil {
        responses.Error(c, http.StatusInternalServerError, err.Error())
        return
    }

    var response []dto.ContactResponse
    for _, u := range contacts {
        response = append(response, dto.ContactResponse{
            ID: u.ID.String(),
            UserID: u.UserID,
            ContactID: u.ContactID,
            Status: u.Status,
        })
    }

    responses.SuccessPaginated(c, "Get data successfully", response, page, limit, total)
}

func (h ContactHandler) FindById(c *gin.Context) {
	id := c.Param("id")
	contactID, err := uuid.Parse(id)
    if err != nil {
        responses.Error(c, http.StatusBadRequest, "Invalid UUID")
        return
    }

	result, errResult := h.uc.FindById(contactID)
	if errResult != nil {
		if ce, ok := errResult.(responses.CodedError); ok {
			responses.Error(c, ce.StatusCode(), ce.Error())
		} else {
			responses.Error(c, http.StatusInternalServerError, errResult.Error())
		}
		return
	}

	response := dto.ContactResponse {
        ID: result.ID.String(),
		UserID: result.UserID,
        ContactID: result.ContactID,
        Status: result.Status,
    }

	responses.Success(c, "Get data successfully", response)
}

func (h *ContactHandler) Truncate(c *gin.Context) {
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
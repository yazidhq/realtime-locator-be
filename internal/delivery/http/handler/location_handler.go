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

type LocationInterface interface {
	Create(location *dto.LocationCreateRequest) (*model.Location, error)
	Update(locationID uuid.UUID, req *dto.LocationUpdateRequest) (*model.Location, error)
	Delete(locationID uuid.UUID) (*model.Location, error)
	FindAll(page, limit int, filters []utils.FilterOptions, sorts []utils.SortOption) ([]model.Location, int, error)
	FindById(locationID uuid.UUID) (*model.Location, error)
	Truncate() (error)
}

type LocationHandler struct {
	uc LocationInterface
}

func NewLocationHandler(uc LocationInterface) *LocationHandler {
	return &LocationHandler{uc: uc}
}

func (h LocationHandler) Create(c *gin.Context) {
	var req dto.LocationCreateRequest
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

func (h *LocationHandler) Update(c *gin.Context) {
    id := c.Param("id")

	locationID, err := uuid.Parse(id)
    if err != nil {
        responses.Error(c, http.StatusBadRequest, "Invalid UUID")
        return
    }

    var req dto.LocationUpdateRequest
    if errReq := c.ShouldBindJSON(&req); errReq != nil {
        responses.Error(c, http.StatusBadRequest, "Invalid request body: " + errReq.Error())
        return
    }

    updated, errUpdate := h.uc.Update(locationID, &req)
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

func (h LocationHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	locationID, err := uuid.Parse(id)
    if err != nil {
        responses.Error(c, http.StatusBadRequest, "Invalid UUID")
        return
    }

	_, errDelete := h.uc.Delete(locationID)
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

func (h LocationHandler) FindAll(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	allowedOps := []string{"=", "like", ">", "<"}
	filters := utils.BuildDynamicFilters(c.Request.URL.Query(), allowedOps)

	allowedFields := []string{"created_at", "latitude", "longitude", "user_id"}
	sorts := utils.BuildDynamicSorts(c.Request.URL.Query(), allowedFields)

	locations, total, err := h.uc.FindAll(page, limit, filters, sorts)
    if err != nil {
        responses.Error(c, http.StatusInternalServerError, err.Error())
        return
    }

    var response []dto.LocationResponse
    for _, u := range locations {
        response = append(response, dto.LocationResponse{
            ID: u.ID.String(),
            UserID: u.UserID,
            Latitude: u.Latitude,
            Longitude: u.Longitude,
        })
    }

    responses.SuccessPaginated(c, "Get data successfully", response, page, limit, total)
}

func (h LocationHandler) FindById(c *gin.Context) {
	id := c.Param("id")
	locationID, err := uuid.Parse(id)
    if err != nil {
        responses.Error(c, http.StatusBadRequest, "Invalid UUID")
        return
    }

	result, errResult := h.uc.FindById(locationID)
	if errResult != nil {
		if ce, ok := errResult.(responses.CodedError); ok {
			responses.Error(c, ce.StatusCode(), ce.Error())
		} else {
			responses.Error(c, http.StatusInternalServerError, errResult.Error())
		}
		return
	}

	response := dto.LocationResponse {
        ID: result.ID.String(),
		UserID: result.UserID,
		Latitude: result.Latitude,
		Longitude: result.Longitude,
    }

	responses.Success(c, "Get data successfully", response)
}

func (h *LocationHandler) Truncate(c *gin.Context) {
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
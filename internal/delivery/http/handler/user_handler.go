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

type UserInterface interface {
	Create(user *dto.UserCreateRequest) (*model.User, error)
	Update(userID uuid.UUID, req *dto.UserUpdateRequest) (*model.User, error)
	Delete(userID uuid.UUID) (*model.User, error)
	FindAll(page, limit int, filters []utils.FilterOptions, sorts []utils.SortOption) ([]model.User, int, error)
	FindById(userID uuid.UUID) (*model.User, error)
	Truncate() (error)
}

type UserHandler struct {
	uc UserInterface
}

func NewUserHandler(uc UserInterface) *UserHandler {
	return &UserHandler{uc: uc}
}

func (h UserHandler) Create(c *gin.Context) {
	var req dto.UserCreateRequest
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

func (h *UserHandler) Update(c *gin.Context) {
    id := c.Param("id")

	userID, err := uuid.Parse(id)
    if err != nil {
        responses.Error(c, http.StatusBadRequest, "Invalid UUID")
        return
    }

    var req dto.UserUpdateRequest
    if errReq := c.ShouldBindJSON(&req); errReq != nil {
        responses.Error(c, http.StatusBadRequest, "Invalid request body: " + errReq.Error())
        return
    }

    updated, errUpdate := h.uc.Update(userID, &req)
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

func (h UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	userID, err := uuid.Parse(id)
    if err != nil {
        responses.Error(c, http.StatusBadRequest, "Invalid UUID")
        return
    }

	_, errDelete := h.uc.Delete(userID)
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

func (h UserHandler) FindAll(c *gin.Context) {
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

	allowedFields := []string{"created_at", "name", "email", "username", "phone_number"}
	sorts := utils.BuildDynamicSorts(c.Request.URL.Query(), allowedFields)

	users, total, err := h.uc.FindAll(page, limit, filters, sorts)
    if err != nil {
        responses.Error(c, http.StatusInternalServerError, err.Error())
        return
    }

    var response []dto.UserResponse
    for _, u := range users {
        response = append(response, dto.UserResponse{
            ID: u.ID.String(),
			Role: string(u.Role),
            Name: u.Name,
            Username: u.Username,
            Email: u.Email,
            PhoneNumber: u.PhoneNumber,
        })
    }

    responses.SuccessPaginated(c, "Get data successfully", response, page, limit, total)
}

func (h UserHandler) FindById(c *gin.Context) {
	id := c.Param("id")
	userID, err := uuid.Parse(id)
    if err != nil {
        responses.Error(c, http.StatusBadRequest, "Invalid UUID")
        return
    }

	result, errResult := h.uc.FindById(userID)
	if errResult != nil {
		if ce, ok := errResult.(responses.CodedError); ok {
			responses.Error(c, ce.StatusCode(), ce.Error())
		} else {
			responses.Error(c, http.StatusInternalServerError, errResult.Error())
		}
		return
	}

	response := dto.UserResponse {
        ID: result.ID.String(),
		Role: string(result.Role),
        Name: result.Name,
        Username: result.Username,
        Email: result.Email,
        PhoneNumber: result.PhoneNumber,
    }

	responses.Success(c, "Get data successfully", response)
}

func (h *UserHandler) Truncate(c *gin.Context) {
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
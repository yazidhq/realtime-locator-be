package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    any 		`json:"data,omitempty"`
}

func Success(c *gin.Context, message string, data any) {
	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func Created(c *gin.Context, message string, data any) {
	c.JSON(http.StatusCreated, APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

type PaginationMeta struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Rows  int `json:"rows"`
	Pages int `json:"pages"`
}

type PaginatedResponse struct {
	Status     string         `json:"status"`
	Message    string         `json:"message"`
	Pagination PaginationMeta `json:"pagination"`
	Data       any    		  `json:"data"`
}

func SuccessPaginated(c *gin.Context, message string, data any, page, limit, total int) {
	pages := (total + limit - 1) / limit
	c.JSON(http.StatusOK, PaginatedResponse{
		Status:  "success",
		Message: message,
		Data:    data,
		Pagination: PaginationMeta{
			Page:       page,
			Limit:      limit,
			Rows:  total,
			Pages: pages,
		},
	})
}

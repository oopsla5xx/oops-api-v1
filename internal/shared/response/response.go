package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/oopsla5xx/oops-api-v1/internal/shared/constants"
	app_errors "github.com/oopsla5xx/oops-api-v1/internal/shared/errors"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Meta struct {
	Page       int `json:"page,omitempty"`
	PerPage    int `json:"per_page,omitempty"`
	Total      int `json:"total,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    data,
	})
}

func Error(c *gin.Context, err error) {
	if appErr, ok := app_errors.IsAppError(err); ok {
		c.JSON(appErr.Status, Response{
			Success: false,
			Error: &ErrorInfo{
				Code:    appErr.Code,
				Message: appErr.Message,
			},
		})
		return
	}

	c.JSON(http.StatusInternalServerError, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    constants.ErrInternalServer,
			Message: "internal server error",
		},
	})
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

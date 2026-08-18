package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/oopsla5xx/oops-api-v1/internal/shared/constants"
	app_errors "github.com/oopsla5xx/oops-api-v1/internal/shared/errors"
)

type Response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  *ErrorBody  `json:"error,omitempty"`
	Meta   *Meta       `json:"meta,omitempty"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Meta struct {
	Page    int `json:"page,omitempty"`
	PerPage int `json:"per_page,omitempty"`
	Total   int `json:"total,omitempty"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Status: "ok",
		Data:   data,
	})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Status: "ok",
		Data:   data,
	})
}

func Error(c *gin.Context, err error) {
	if appErr, ok := app_errors.IsAppError(err); ok {
		c.JSON(appErr.StatusCode, Response{
			Status: "error",
			Error: &ErrorBody{
				Code:    appErr.Code,
				Message: appErr.Message,
			},
		})
		return
	}

	c.JSON(http.StatusInternalServerError, Response{
		Status: "error",
		Error: &ErrorBody{
			Code:    constants.ErrInternalServer,
			Message: "internal server error",
		},
	})
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

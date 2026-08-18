package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/oopsla5xx/oops-api-v1/internal/shared/constants"
	"github.com/oopsla5xx/oops-api-v1/internal/shared/response"
)

func Recovery(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Error("panic recovered",
					zap.Any("error", err),
					zap.String("request_id", c.GetString(constants.HeaderRequestID)),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
				)
				c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{
					Status: "error",
					Error: &response.ErrorBody{
						Code:    constants.ErrInternalServer,
						Message: "internal server error",
					},
				})
			}
		}()
		c.Next()
	}
}

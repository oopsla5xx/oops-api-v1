package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/oopsla5xx/oops-api-v1/internal/shared/constants"
	"github.com/oopsla5xx/oops-api-v1/internal/shared/response"
)

// Timeout sets a deadline on the request context.
// Downstream handlers must respect ctx.Done() for this to take effect.
func Timeout(duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), duration)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
		c.Next()

		if ctx.Err() == context.DeadlineExceeded && !c.Writer.Written() {
			c.AbortWithStatusJSON(http.StatusGatewayTimeout, response.Response{
				Status: "error",
				Error: &response.ErrorBody{
					Code:    constants.ErrGatewayTimeout,
					Message: "request timeout",
				},
			})
		}
	}
}

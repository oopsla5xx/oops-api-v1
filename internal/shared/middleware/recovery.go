package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
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
				noticeError(c, err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{
					Success: false,
					Error: &response.ErrorInfo{
						Code:    constants.ErrInternalServer,
						Message: "internal server error",
					},
				})
			}
		}()
		c.Next()
	}
}

// noticeError reports a recovered panic to the New Relic transaction, if
// one is present in the request. nrgin.Transaction returns nil when New
// Relic is disabled or nrgin.Middleware did not run, and *Transaction
// methods are nil-safe, so this is a no-op in that case.
func noticeError(c *gin.Context, recovered any) {
	txn := nrgin.Transaction(c)
	if err, ok := recovered.(error); ok {
		txn.NoticeError(err)
		return
	}
	txn.NoticeError(fmt.Errorf("panic: %v", recovered))
}

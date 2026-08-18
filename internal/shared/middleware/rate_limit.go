package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"github.com/oopsla5xx/oops-api-v1/internal/shared/constants"
	"github.com/oopsla5xx/oops-api-v1/internal/shared/response"
)

// RateLimit applies a per-IP token-bucket rate limiter.
// rps is the sustained request rate; burst is the maximum burst size.
//
// Note: limiters are stored in an in-memory sync.Map and never evicted.
// This is acceptable for a single-replica deployment. For multi-replica
// production use, replace with a Redis-based sliding-window implementation.
func RateLimit(rps float64, burst int) gin.HandlerFunc {
	var limiters sync.Map

	return func(c *gin.Context) {
		ip := c.ClientIP()

		v, _ := limiters.LoadOrStore(ip, rate.NewLimiter(rate.Limit(rps), burst))
		limiter := v.(*rate.Limiter)

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, response.Response{
				Status: "error",
				Error: &response.ErrorBody{
					Code:    constants.ErrTooManyRequests,
					Message: "too many requests",
				},
			})
			return
		}
		c.Next()
	}
}

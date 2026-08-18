package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/oopsla5xx/oops-api-v1/internal/shared/constants"
)

func CORS(allowedOrigins []string) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:  allowedOrigins,
		AllowMethods:  []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Accept", "Authorization", constants.HeaderRequestID},
		ExposeHeaders: []string{constants.HeaderRequestID},
		MaxAge:        12 * time.Hour,
	})
}

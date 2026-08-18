package app

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/oopsla5xx/oops-api-v1/docs"
	health_module "github.com/oopsla5xx/oops-api-v1/internal/modules/health"
	"github.com/oopsla5xx/oops-api-v1/internal/shared/constants"
	"github.com/oopsla5xx/oops-api-v1/internal/shared/middleware"
)

func newRouter(deps *dependencies) *gin.Engine {
	if deps.cfg.App.Env == constants.Production {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.Use(middleware.Recovery(deps.log))
	r.Use(middleware.RequestID())
	r.Use(middleware.CORS(deps.cfg.CORS.AllowedOrigins))
	r.Use(middleware.RequestLogger(deps.log))
	r.Use(middleware.Timeout(deps.cfg.Server.ReadTimeout))
	r.Use(middleware.RateLimit(deps.cfg.RateLimit.RequestsPerSecond, deps.cfg.RateLimit.Burst))

	if deps.cfg.App.Env != constants.Production {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	v1 := r.Group(constants.APIVersionV1)
	{
		health_module.New().Register(v1)
	}

	return r
}

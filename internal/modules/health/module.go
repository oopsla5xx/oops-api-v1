package health

import (
	health_query "github.com/oopsla5xx/oops-api-v1/internal/modules/health/application/query"
	health_handler "github.com/oopsla5xx/oops-api-v1/internal/modules/health/interface"
)

func New() *health_handler.Handler {
	return health_handler.NewHandler(health_query.NewHealthQuery())
}

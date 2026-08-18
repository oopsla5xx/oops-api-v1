package query

import (
	"github.com/oopsla5xx/oops-api-v1/internal/shared/constants"
	"github.com/oopsla5xx/oops-api-v1/internal/shared/version"
)

type HealthResult struct {
	Status  string
	Service string
	Version string
}

type HealthQuery struct{}

func NewHealthQuery() *HealthQuery {
	return &HealthQuery{}
}

func (q *HealthQuery) Execute() HealthResult {
	return HealthResult{
		Status:  "ok",
		Service: constants.AppName,
		Version: version.Version,
	}
}

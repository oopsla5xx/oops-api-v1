package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/oopsla5xx/oops-api-v1/internal/modules/health/application/query"
)

// Handler handles HTTP requests for the health module.
type Handler struct {
	query *query.HealthQuery
}

func NewHandler(q *query.HealthQuery) *Handler {
	return &Handler{query: q}
}

// Register wires the handler routes onto the given router group.
func (h *Handler) Register(rg *gin.RouterGroup) {
	rg.GET("/health", h.Health)
}

// Health godoc
// @Summary      Health check
// @Description  Returns service liveness status. Does NOT check DB/Redis — use this for liveness probes only.
// @Tags         health
// @Produce      json
// @Success      200  {object}  healthResponse
// @Router       /health [get]
//
// Uses a flat response (not response.OK wrapper) intentionally — liveness probes
// expect a stable, minimal shape that orchestrators (k8s, ECS) can parse cheaply.
func (h *Handler) Health(c *gin.Context) {
	result := h.query.Execute()
	c.JSON(http.StatusOK, healthResponse{
		Status:  result.Status,
		Service: result.Service,
		Version: result.Version,
	})
}

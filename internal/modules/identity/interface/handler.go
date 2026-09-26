package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/oopsla5xx/oops-api-v1/internal/modules/identity/application/command"
	"github.com/oopsla5xx/oops-api-v1/internal/shared/constants"
	app_errors "github.com/oopsla5xx/oops-api-v1/internal/shared/errors"
	"github.com/oopsla5xx/oops-api-v1/internal/shared/response"
)

type Handler struct {
	createUser *command.CreateUserCommand
}

func NewHandler(createUser *command.CreateUserCommand) *Handler {
	return &Handler{createUser: createUser}
}

func (h *Handler) Register(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	auth.POST("/register", h.CreateUser)
}

// CreateUser godoc
// @Summary      Register a user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      createUserRequest  true  "user data"
// @Success      201      {object}  userResponse
// @Failure      400      {object}  response.Response
// @Failure      409      {object}  response.Response
// @Router       /auth/register [post]
func (h *Handler) CreateUser(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, app_errors.Wrap(err, constants.ErrValidation, err.Error(), http.StatusBadRequest))
		return
	}

	if validationErr := validateRequest(req); validationErr != nil {
		response.Error(c, validationErr)
		return
	}

	user, err := h.createUser.Execute(c.Request.Context(), command.CreateUserInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		response.Error(c, mapError(err))
		return
	}

	response.Created(c, toResponse(user))
}

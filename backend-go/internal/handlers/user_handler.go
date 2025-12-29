package handlers

import (
	"quiz-game-backend/internal/database"
	"quiz-game-backend/internal/dto"
	"quiz-game-backend/internal/repository"
	"quiz-game-backend/internal/services"
	"quiz-game-backend/internal/utils"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandler struct {
	service   services.UserService
	validator *validator.Validate
}

func NewUserHandler(db *database.Database) *UserHandler {
	repo := repository.NewUserResponse(db.DB)
	service := services.NewUserService(repo)
	return &UserHandler{
		service:   service,
		validator: validator.New(),
	}
}

// GetUser godoc
// @Summary Get user by ID
// @Tags users
// @Param id path string true "User ID"
// @Success 200 {object} dto.UserResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.BadRequest(c, "Invalid user ID")
	}
	user, err := h.service.GetUserByID(c.Context(), id)
	if err != nil {
		return utils.NotFound(c, err.Error())
	}
	return utils.Success(c, dto.ToUserResponse(user))
}

// UpdateProfile godoc
// @Summary Update user profile
// @Tags users
// @Param id path string true "User ID"
// @Param body body dto.UpdateProfileRequest true "Profile data"
// @Success 200 {object} dto.UserResponse
// @Router /users/{id} [put]
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.BadRequest(c, "Invalid user ID")
	}
	var req dto.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}
	if err := h.validator.Struct(req); err != nil {
		return utils.BadRequest(c, err.Error())
	}
	user, err := h.service.UpdateProfile(c.Context(), id, req)
	if err != nil {
		return utils.InternalServerError(c, err.Error())
	}
	return utils.Success(c, dto.ToUserResponse(user))
}

// ListUsers godoc
// @Summary List users with pagination
// @Tags users
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} dto.UserListResponse
// @Router /users [get]
func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	result, err := h.service.ListUsers(c.Context(), page, limit)
	if err != nil {
		return utils.InternalServerError(c, err.Error())
	}
	return utils.Paginated(c, result.Users, page, limit, result.Total)
}

package handlers

import (
	"quiz-game-backend/internal/cache"
	"quiz-game-backend/internal/config"
	"quiz-game-backend/internal/database"
	"quiz-game-backend/internal/dto"
	"quiz-game-backend/internal/repository"
	"quiz-game-backend/internal/services"
	"quiz-game-backend/internal/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AuthHandler struct {
	service     services.AuthService
	userService services.UserService // เพิ่ม
	validator   *validator.Validate
	cfg         *config.Config
	userCache   *cache.UserCache
}

func NewAuthHandler(db *database.Database, cfg *config.Config, redisClient *cache.RedisClient) *AuthHandler {
	repo := repository.NewUserResponse(db.DB)
	authService := services.NewAuthService(repo, cfg)
	userService := services.NewUserService(repo) // เพิ่ม
	return &AuthHandler{
		service:     authService,
		userService: userService, // เพิ่ม
		validator:   validator.New(),
		cfg:         cfg,
		userCache:   cache.NewUserCache(redisClient),
	}
}

// Register สมัครสมาชิก
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	if err := h.validator.Struct(req); err != nil {
		return utils.BadRequest(c, err.Error())
	}
	response, err := h.service.Register(c.Context(), req)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}
	// ตั้งค่า cookie
	utils.SetAuthCookie(c, response.Token, h.cfg.JWT.Expiration, h.cfg.IsProduction())
	return utils.Created(c, response.User, "Registration successful")
}

// Login เข้าสู่ระบบ
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}
	if err := h.validator.Struct(req); err != nil {
		return utils.BadRequest(c, err.Error())
	}
	response, err := h.service.Login(c.Context(), req)
	if err != nil {
		return utils.Unauthorized(c, err.Error())
	}
	// ตั้งค่า cookie
	utils.SetAuthCookie(c, response.Token, h.cfg.JWT.Expiration, h.cfg.IsProduction())
	// set cache
	h.userCache.Set(c.Context(), &response.User)
	
	return utils.Success(c, response.User, "Login successful")
}

// Logout ออกจากระบบ (ต้องผ่าน auth middleware)
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// ดึง userID จาก context (ต้องผ่าน auth middleware)
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if ok {
		// ลบ cache ใน Redis
		h.userCache.Delete(c.Context(), userID.String())
	}

	// ลบ cookie
	utils.ClearAuthCookie(c)
	return utils.Success(c, nil, "Logout successful")
}

// GetMe - ดึงจาก cache ก่อน
func (h *AuthHandler) GetMe(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	// 1. ดึงจาก cache ก่อน
	cached, err := h.userCache.Get(c.Context(), userID.String())
	if err == nil && cached != nil {
		return utils.Success(c, cached) // Cache hit!
	}
	// 2. Cache miss - ดึงจาก DB
	user, err := h.userService.GetUserByID(c.Context(), userID) // แก้จาก h.userService
	if err != nil {
		return utils.NotFound(c, "User not found")
	}

	userResp := dto.ToUserResponse(user)
	return utils.Success(c, userResp)
}

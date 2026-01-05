package handlers

import (
	"log"
	"quiz-game-backend/internal/dto"
	"quiz-game-backend/internal/services"
	"quiz-game-backend/internal/utils"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type QuizHandler struct {
	service   services.QuizService
	validator *validator.Validate
}

func NewQuizHandler(service services.QuizService) *QuizHandler {
	return &QuizHandler{
		service:   service,
		validator: validator.New(),
	}
}

// CreateQuiz godoc
// @Summary Create a new quiz
// @Tags quizzes
// @Security BearerAuth
// @Param body body dto.CreateQuizRequest true "Quiz data"
// @Success 201 {object} dto.QuizResponse
// @Router /quizzes [post]
func (h *QuizHandler) CreateQuiz(c *fiber.Ctx) error {
	var req dto.CreateQuizRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}
	if err := h.validator.Struct(req); err != nil {
		return utils.BadRequest(c, err.Error())
	}
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "User not authenticated")
	}

	quiz, err := h.service.CreateQuiz(c.Context(), userID, req)
	if err != nil {
		return utils.InternalServerError(c, err.Error())
	}
	return utils.Created(c, quiz)
}

// GetQuizByID godoc
// @Summary Get quiz by ID
// @Tags quizzes
// @Param id path string true "Quiz ID"
// @Success 200 {object} dto.QuizResponse
// @Router /quizzes/{id} [get]
func (h *QuizHandler) GetQuizByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.BadRequest(c, "Invalid quiz ID")
	}

	quiz, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return utils.NotFound(c, "Quiz not found")
	}
	return utils.Success(c, quiz)
}

// ListQuizzes godoc
// @Summary List all public quizzes
// @Tags quizzes
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param category query string false "Filter by category"
// @Success 200 {object} dto.QuizListResponse
// @Router /quizzes [get]
func (h *QuizHandler) ListQuizzes(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	category := c.Query("category", "")

	result, err := h.service.ListQuizzes(c.Context(), page, limit, category)
	if err != nil {
		return utils.InternalServerError(c, err.Error())
	}
	return utils.Paginated(c, result.Quizzes, page, limit, int(result.Total))
}

// UpdateQuiz godoc
// @Summary Update quiz (owner only)
// @Tags quizzes
// @Security BearerAuth
// @Param id path string true "Quiz ID"
// @Param body body dto.UpdateQuizRequest true "Quiz data"
// @Success 200 {object} dto.QuizResponse
// @Router /quizzes/{id} [put]
func (h *QuizHandler) UpdateQuiz(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.BadRequest(c, "Invalid quiz ID")
	}

	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "User not authenticated")
	}

	var req dto.UpdateQuizRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}
	if err := h.validator.Struct(req); err != nil {
		return utils.BadRequest(c, err.Error())
	}

	quiz, err := h.service.UpdateQuiz(c.Context(), id, userID, req)
	if err != nil {
		return utils.InternalServerError(c, err.Error())
	}
	return utils.Success(c, quiz)
}

// DeleteQuiz godoc
// @Summary Delete quiz (owner only)
// @Tags quizzes
// @Security BearerAuth
// @Param id path string true "Quiz ID"
// @Success 200 {object} map[string]string
// @Router /quizzes/{id} [delete]
func (h *QuizHandler) DeleteQuiz(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.BadRequest(c, "Invalid quiz ID")
	}

	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "User not authenticated")
	}

	if err := h.service.DeleteQuiz(c.Context(), id, userID); err != nil {
		return utils.InternalServerError(c, err.Error())
	}
	return utils.Success(c, fiber.Map{"message": "Quiz deleted successfully"})
}

// GetMyQuizzes godoc
// @Summary Get current user's quizzes
// @Tags quizzes
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} dto.QuizListResponse
// @Router /quizzes/my [get]
func (h *QuizHandler) GetMyQuizzes(c *fiber.Ctx) error {
	log.Println("รวย")
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "User not authenticated")
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	result, err := h.service.GetByCreator(c.Context(), userID, page, limit)
	if err != nil {
		return utils.InternalServerError(c, err.Error())
	}
	return utils.Paginated(c, result.Quizzes, page, limit, int(result.Total))
}

// AddQuestion godoc
// @Summary Add question to quiz
// @Tags quizzes
// @Security BearerAuth
// @Param id path string true "Quiz ID"
// @Param body body dto.CreateQuestionRequest true "Question data"
// @Success 201 {object} dto.QuestionResponse
// @Router /quizzes/{id}/questions [post]
func (h *QuizHandler) AddQuestion(c *fiber.Ctx) error {
	quizID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.BadRequest(c, "Invalid quiz ID")
	}

	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "User not authenticated")
	}

	var req dto.CreateQuestionRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}
	if err := h.validator.Struct(req); err != nil {
		return utils.BadRequest(c, err.Error())
	}

	question, err := h.service.AddQuestion(c.Context(), quizID, userID, req)
	if err != nil {
		return utils.InternalServerError(c, err.Error())
	}
	return utils.Created(c, question)
}

// UpdateQuestion godoc
// @Summary Update question
// @Tags quizzes
// @Security BearerAuth
// @Param id path string true "Quiz ID"
// @Param questionId path string true "Question ID"
// @Param body body dto.CreateQuestionRequest true "Question data"
// @Success 200 {object} dto.QuestionResponse
// @Router /quizzes/{id}/questions/{questionId} [put]
func (h *QuizHandler) UpdateQuestion(c *fiber.Ctx) error {
	questionID, err := uuid.Parse(c.Params("questionId"))
	if err != nil {
		return utils.BadRequest(c, "Invalid question ID")
	}

	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "User not authenticated")
	}

	var req dto.CreateQuestionRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}
	if err := h.validator.Struct(req); err != nil {
		return utils.BadRequest(c, err.Error())
	}

	question, err := h.service.UpdateQuestion(c.Context(), questionID, userID, req)
	if err != nil {
		return utils.InternalServerError(c, err.Error())
	}
	return utils.Success(c, question)
}

// DeleteQuestion godoc
// @Summary Delete question
// @Tags quizzes
// @Security BearerAuth
// @Param id path string true "Quiz ID"
// @Param questionId path string true "Question ID"
// @Success 200 {object} map[string]string
// @Router /quizzes/{id}/questions/{questionId} [delete]
func (h *QuizHandler) DeleteQuestion(c *fiber.Ctx) error {
	questionID, err := uuid.Parse(c.Params("questionId"))
	if err != nil {
		return utils.BadRequest(c, "Invalid question ID")
	}

	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "User not authenticated")
	}

	if err := h.service.DeleteQuestion(c.Context(), questionID, userID); err != nil {
		return utils.InternalServerError(c, err.Error())
	}
	return utils.Success(c, fiber.Map{"message": "Question deleted successfully"})
}

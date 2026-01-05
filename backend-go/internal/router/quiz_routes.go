package router

import (
	"quiz-game-backend/internal/config"
	"quiz-game-backend/internal/database"
	"quiz-game-backend/internal/handlers"
	"quiz-game-backend/internal/middleware"
	"quiz-game-backend/internal/repository"
	"quiz-game-backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

func SetQuizRoutes(api fiber.Router, db *database.Database, cfg *config.Config) {
	// Dependency Injection: Repository -> Service -> Handler
	repo := repository.NewQuizRepository(db.DB)
	service := services.NewQuizService(repo)
	handler := handlers.NewQuizHandler(service)

	quizzes := api.Group("/quizzes")

	// Public routes (ไม่ต้อง login)
	quizzes.Get("/", handler.ListQuizzes) // List all public quizzes

	// Protected routes (ต้อง login)
	protected := quizzes.Group("", middleware.Auth(cfg))

	// ⚠️ IMPORTANT: Specific routes MUST come BEFORE parameterized routes!
	// /my ต้องอยู่ก่อน /:id ไม่งั้น "my" จะถูกจับเป็น :id
	protected.Get("/my", handler.GetMyQuizzes) // Get logged-in user's quizzes

	// Now the parameterized routes
	quizzes.Get("/:id", handler.GetQuizByID)     // Get quiz details (public)
	protected.Post("/", handler.CreateQuiz)      // Create new quiz
	protected.Put("/:id", handler.UpdateQuiz)    // Update quiz (owner only)
	protected.Delete("/:id", handler.DeleteQuiz) // Delete quiz (owner only)

	// Question routes (nested under quiz, protected)
	protected.Post("/:id/questions", handler.AddQuestion) // Add question to quiz
	protected.Put("/:id/questions/:questionId", handler.UpdateQuestion)
	protected.Delete("/:id/questions/:questionId", handler.DeleteQuestion)
}

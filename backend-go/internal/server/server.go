package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"quiz-game-backend/internal/cache"
	"quiz-game-backend/internal/config"
	"quiz-game-backend/internal/database"
	"quiz-game-backend/internal/router"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/google/uuid"
)

type Server struct {
	app         *fiber.App
	config      *config.Config
	db          *database.Database
	redisClient *cache.RedisClient
}

func New(cfg *config.Config, db *database.Database, redisClient *cache.RedisClient) *Server {
	app := fiber.New(fiber.Config{
		AppName:               "Quiz Game API",
		ServerHeader:          "Quiz-Game",
		StrictRouting:         true,
		CaseSensitive:         true,
		ReadTimeout:           cfg.Server.ReadTimeout,
		WriteTimeout:          cfg.Server.WriteTimeout,
		IdleTimeout:           time.Minute,
		ErrorHandler:          customErrorHandler,
		DisableStartupMessage: true,
	})

	// Security middleware (ป้องกัน common attacks)
	app.Use(helmet.New(helmet.Config{
		XSSProtection:           "1; mode=block",
		ContentTypeNosniff:      "nosniff",
		XFrameOptions:           "DENY",
		ReferrerPolicy:          "no-referrer",
		CrossOriginOpenerPolicy: "same-origin",
	}))

	// Request ID (ติดตาม request แต่ละตัว)
	app.Use(requestid.New(requestid.Config{
		Generator: func() string {
			return uuid.New().String()
		},
	}))

	// Recover from panics
	app.Use(recover.New(recover.Config{
		EnableStackTrace: cfg.IsDevelopment(),
	}))

	// Logging
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${locals:requestid} ${status} - ${method} ${path} (${latency})\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Asia/Bangkok",
	}))

	// CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     getAllowedOrigins(cfg),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length",
		MaxAge:           3600,
	}))

	// Compression (gzip response)
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	// Rate limiting (ป้องกัน DDoS)
	if cfg.IsProduction() {
		app.Use(limiter.New(limiter.Config{
			Max:        100,             // 100 requests
			Expiration: 1 * time.Minute, // per minute
			KeyGenerator: func(c *fiber.Ctx) string {
				return c.IP()
			},
			LimitReached: func(c *fiber.Ctx) error {
				return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
					"error": "Rate limit exceeded",
				})
			},
		}))
	}

	return &Server{
		app:         app,
		config:      cfg,
		db:          db,
		redisClient: redisClient,
	}

}

func (s *Server) Run() error {
	// Setup routes
	router.Setup(s.app, s.config, s.db, s.redisClient)

	// Handle graceful shutdown
	return s.runWithGracefulShutdown()
}

func (s *Server) Shutdown() error {
	// Close database
	if err := s.db.Close(); err != nil {
		log.Printf("Failed to close database: %v", err)
	}
	return s.app.Shutdown()
}

func (s *Server) runWithGracefulShutdown() error {
	// Start server in goroutine
	addr := fmt.Sprintf("%s:%s", s.config.Server.Address, s.config.Server.Port)

	go func() {
		log.Printf("🚀 Server starting on %s", addr)
		log.Printf("📊 Environment: %s", s.config.Server.Environment)
		if err := s.app.Listen(addr); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()
	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Println("🛑 Shutting down server...")
	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), s.config.Server.ShutdownTimeout)
	defer cancel()
	if err := s.app.ShutdownWithContext(ctx); err != nil {
		log.Printf("❌ Server forced to shutdown: %v", err)
		return err
	}
	log.Println("✅ Server stopped gracefully")

	return nil
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}
	// Log error
	log.Printf("Error: %v (RequestID: %s)", err, c.Locals("requestid"))
	return c.Status(code).JSON(fiber.Map{
		"error":      message,
		"request_id": c.Locals("requestid"),
	})
}

func getAllowedOrigins(cfg *config.Config) string {
	if cfg.IsDevelopment() {
		return "http://localhost:3000"
	}
	// Production - เพิ่ม domain จริง
	return "https://your-frontend-domain.com,https://your-app.vercel.app"
}

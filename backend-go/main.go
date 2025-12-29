package main

import (
	"log"
	"quiz-game-backend/internal/cache"
	"quiz-game-backend/internal/config"
	"quiz-game-backend/internal/database"
	"quiz-game-backend/internal/server"
)

func main() {
	// load cfg
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	// connect db
	db, err := database.NewPostgres(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}
	defer db.Close()

	// Auto migrate
	if err := db.AutoMigrate(); err != nil {
		log.Fatal("Failed to migrate:", err)
	}

	// connet redis
	redisClient, err := cache.NewRedisClient(cfg.Redis.URL, cfg.Redis.DB)
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	defer redisClient.Close()
	log.Println("✅ Redis connected")

	// Log startup info
	log.Printf("Starting Quiz Game API")
	log.Printf("Environment: %s", cfg.Server.Environment)
	log.Printf("Port: %s", cfg.Server.Port)

	// Create and run server
	srv := server.New(cfg, db, redisClient)
	if err := srv.Run(); err != nil {
		log.Fatal("Server error:", err)
	}

}

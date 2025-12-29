package database

import (
	"fmt"
	"log"
	"quiz-game-backend/internal/config"
	"quiz-game-backend/internal/models"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgres(cfg *config.DatabaseConfig) (*Database, error) {
	var db *gorm.DB
	var err error

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	// Retry connection
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(cfg.URL), gormConfig)
		if err == nil {
			break
		}

		log.Printf("failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(time.Duration(i+1) * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect after %d attempts: %w", maxRetries, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	log.Println("✅ Database connected successfully")
	return New(db), nil

}

func (d *Database) AutoMigrate() error {
	err := d.DB.AutoMigrate(
		&models.User{},
		&models.Quiz{},
		&models.Question{},
		&models.GameRoom{},
		&models.GameParticipant{},
		&models.GameAnswer{},
		&models.Message{},
	)

	if err != nil {
		return err
	}

	log.Println("✅ Database migrated successfully")
	return nil
}

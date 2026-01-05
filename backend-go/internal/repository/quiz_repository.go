package repository

import (
	"context"
	"log"
	"quiz-game-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuizRepository interface {
	CreateQuiz(ctx context.Context, quiz models.Quiz) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Quiz, error)
	FindByIDWithQuestions(ctx context.Context, id uuid.UUID) (*models.Quiz, error)
	UpdateQuiz(ctx context.Context, quiz models.Quiz) error
	DeleteQuiz(ctx context.Context, id uuid.UUID) error
	ListQuiz(ctx context.Context, limit, offset int) ([]models.Quiz, int64, error)
	ListByCategory(ctx context.Context, category string, limit, offset int) ([]models.Quiz, int64, error)
	ListByCreator(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.Quiz, int64, error)
	CreateQuestion(ctx context.Context, question models.Question) error
	UpdateQuestion(ctx context.Context, question models.Question) error
	DeleteQuestion(ctx context.Context, id uuid.UUID) error
}

type quizRepository struct {
	db *gorm.DB
}

func NewQuizRepository(db *gorm.DB) QuizRepository {
	return &quizRepository{
		db: db,
	}
}

func (r *quizRepository) CreateQuiz(ctx context.Context, quiz models.Quiz) error {
	log.Println("CreateQuiz", quiz)
	return r.db.WithContext(ctx).Create(&quiz).Error
}

func (r *quizRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Quiz, error) {
	var quiz models.Quiz
	if err := r.db.WithContext(ctx).First(&quiz, id).Error; err != nil {
		return nil, err
	}
	return &quiz, nil
}

func (r *quizRepository) FindByIDWithQuestions(ctx context.Context, id uuid.UUID) (*models.Quiz, error) {
	var quiz models.Quiz
	if err := r.db.WithContext(ctx).Preload("Questions").First(&quiz, id).Error; err != nil {
		return nil, err
	}
	return &quiz, nil
}

func (r *quizRepository) UpdateQuiz(ctx context.Context, quiz models.Quiz) error {
	return r.db.WithContext(ctx).Save(&quiz).Error
}

func (r *quizRepository) DeleteQuiz(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Quiz{}, "id = ?", id).Error
}

func (r *quizRepository) ListQuiz(ctx context.Context, limit, offset int) ([]models.Quiz, int64, error) {
	var quizzes []models.Quiz
	var total int64
	if err := r.db.WithContext(ctx).Model(&models.Quiz{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&quizzes).Error
	return quizzes, total, err
}

func (r *quizRepository) ListByCategory(ctx context.Context, category string, limit, offset int) ([]models.Quiz, int64, error) {
	var quizzes []models.Quiz
	var total int64
	if err := r.db.WithContext(ctx).Model(&models.Quiz{}).Where("category = ?", category).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.WithContext(ctx).
		Where("category = ?", category).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&quizzes).Error
	return quizzes, total, err
}

func (r *quizRepository) ListByCreator(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.Quiz, int64, error) {
	var quizzes []models.Quiz
	var total int64
	if err := r.db.WithContext(ctx).Model(&models.Quiz{}).Where("created_by_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.WithContext(ctx).
		Where("created_by_id = ?", userID).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&quizzes).Error
	return quizzes, total, err
}

func (r *quizRepository) CreateQuestion(ctx context.Context, question models.Question) error {
	return r.db.WithContext(ctx).Create(&question).Error
}

func (r *quizRepository) UpdateQuestion(ctx context.Context, question models.Question) error {
	return r.db.WithContext(ctx).Save(&question).Error
}

func (r *quizRepository) DeleteQuestion(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Question{}, "id = ?", id).Error
}

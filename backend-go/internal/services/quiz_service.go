package services

import (
	"context"
	"errors"
	"quiz-game-backend/internal/dto"
	"quiz-game-backend/internal/models"
	"quiz-game-backend/internal/repository"

	"github.com/google/uuid"
)

type QuizService interface {
	// Quiz CRUD
	CreateQuiz(ctx context.Context, userID uuid.UUID, req dto.CreateQuizRequest) (*dto.QuizResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*dto.QuizResponse, error)
	UpdateQuiz(ctx context.Context, id, userID uuid.UUID, req dto.UpdateQuizRequest) (*dto.QuizResponse, error)
	DeleteQuiz(ctx context.Context, id, userID uuid.UUID) error
	ListQuizzes(ctx context.Context, page, limit int, category string) (*dto.QuizListResponse, error)
	GetByCreator(ctx context.Context, userID uuid.UUID, page, limit int) (*dto.QuizListResponse, error)

	// Question CRUD
	AddQuestion(ctx context.Context, quizID, userID uuid.UUID, req dto.CreateQuestionRequest) (*dto.QuestionResponse, error)
	UpdateQuestion(ctx context.Context, questionID, userID uuid.UUID, req dto.CreateQuestionRequest) (*dto.QuestionResponse, error)
	DeleteQuestion(ctx context.Context, questionID, userID uuid.UUID) error
}

type quizService struct {
	repo repository.QuizRepository
}

func NewQuizService(repo repository.QuizRepository) QuizService {
	return &quizService{repo: repo}
}

// CreateQuiz creates a new quiz
func (s *quizService) CreateQuiz(ctx context.Context, userID uuid.UUID, req dto.CreateQuizRequest) (*dto.QuizResponse, error) {
	quiz := models.Quiz{
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
		Difficulty:  req.Difficulty,
		TimeLimit:   req.TimeLimit,
		CreatedByID: userID,
	}

	if err := s.repo.CreateQuiz(ctx, quiz); err != nil {
		return nil, err
	}

	return &dto.QuizResponse{
		ID:          quiz.ID.String(),
		Title:       quiz.Title,
		Description: quiz.Description,
		Category:    quiz.Category,
		Difficulty:  quiz.Difficulty,
		TimeLimit:   quiz.TimeLimit,
		IsPublished: quiz.IsPublished,
		CreatedByID: quiz.CreatedByID.String(),
	}, nil
}

// GetByID gets a quiz by ID with questions
func (s *quizService) GetByID(ctx context.Context, id uuid.UUID) (*dto.QuizResponse, error) {
	quiz, err := s.repo.FindByIDWithQuestions(ctx, id)
	if err != nil {
		return nil, err
	}

	response := dto.QuizResponse{
		ID:          quiz.ID.String(),
		Title:       quiz.Title,
		Description: quiz.Description,
		Category:    quiz.Category,
		Difficulty:  quiz.Difficulty,
		TimeLimit:   quiz.TimeLimit,
		IsPublished: quiz.IsPublished,
		CreatedByID: quiz.CreatedByID.String(),
	}

	// Convert questions
	for _, q := range quiz.Questions {
		response.Questions = append(response.Questions, dto.QuestionResponse{
			ID:           q.ID.String(),
			QuestionText: q.QuestionText,
			OptionA:      q.OptionA,
			OptionB:      q.OptionB,
			OptionC:      q.OptionC,
			OptionD:      q.OptionD,
			Points:       q.Points,
			OrderIndex:   q.OrderIndex,
		})
	}

	return &response, nil
}

// UpdateQuiz updates a quiz (owner only)
func (s *quizService) UpdateQuiz(ctx context.Context, id, userID uuid.UUID, req dto.UpdateQuizRequest) (*dto.QuizResponse, error) {
	quiz, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check ownership
	if quiz.CreatedByID != userID {
		return nil, errors.New("you don't have permission to update this quiz")
	}

	// Update fields
	if req.Title != "" {
		quiz.Title = req.Title
	}
	if req.Description != "" {
		quiz.Description = req.Description
	}
	if req.Category != "" {
		quiz.Category = req.Category
	}
	if req.Difficulty != "" {
		quiz.Difficulty = req.Difficulty
	}
	if req.TimeLimit > 0 {
		quiz.TimeLimit = req.TimeLimit
	}
	quiz.IsPublished = req.IsPublished

	if err := s.repo.UpdateQuiz(ctx, *quiz); err != nil {
		return nil, err
	}

	return &dto.QuizResponse{
		ID:          quiz.ID.String(),
		Title:       quiz.Title,
		Description: quiz.Description,
		Category:    quiz.Category,
		Difficulty:  quiz.Difficulty,
		TimeLimit:   quiz.TimeLimit,
		IsPublished: quiz.IsPublished,
		CreatedByID: quiz.CreatedByID.String(),
	}, nil
}

// DeleteQuiz deletes a quiz (owner only)
func (s *quizService) DeleteQuiz(ctx context.Context, id, userID uuid.UUID) error {
	quiz, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Check ownership
	if quiz.CreatedByID != userID {
		return errors.New("you don't have permission to delete this quiz")
	}

	return s.repo.DeleteQuiz(ctx, id)
}

// ListQuizzes lists all public quizzes with optional category filter
func (s *quizService) ListQuizzes(ctx context.Context, page, limit int, category string) (*dto.QuizListResponse, error) {
	offset := (page - 1) * limit

	var quizzes []models.Quiz
	var total int64
	var err error

	if category != "" {
		quizzes, total, err = s.repo.ListByCategory(ctx, category, limit, offset)
	} else {
		quizzes, total, err = s.repo.ListQuiz(ctx, limit, offset)
	}

	if err != nil {
		return nil, err
	}

	var response []dto.QuizResponse
	for _, quiz := range quizzes {
		response = append(response, dto.QuizResponse{
			ID:          quiz.ID.String(),
			Title:       quiz.Title,
			Description: quiz.Description,
			Category:    quiz.Category,
			Difficulty:  quiz.Difficulty,
			TimeLimit:   quiz.TimeLimit,
			IsPublished: quiz.IsPublished,
			CreatedByID: quiz.CreatedByID.String(),
		})
	}

	return &dto.QuizListResponse{
		Quizzes: response,
		Total:   total,
	}, nil
}

// GetByCreator gets quizzes by creator (user)
func (s *quizService) GetByCreator(ctx context.Context, userID uuid.UUID, page, limit int) (*dto.QuizListResponse, error) {
	offset := (page - 1) * limit

	quizzes, total, err := s.repo.ListByCreator(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	var response []dto.QuizResponse
	for _, quiz := range quizzes {
		response = append(response, dto.QuizResponse{
			ID:          quiz.ID.String(),
			Title:       quiz.Title,
			Description: quiz.Description,
			Category:    quiz.Category,
			Difficulty:  quiz.Difficulty,
			TimeLimit:   quiz.TimeLimit,
			IsPublished: quiz.IsPublished,
			CreatedByID: quiz.CreatedByID.String(),
		})
	}

	return &dto.QuizListResponse{
		Quizzes: response,
		Total:   total,
	}, nil
}

// AddQuestion adds a question to a quiz
func (s *quizService) AddQuestion(ctx context.Context, quizID, userID uuid.UUID, req dto.CreateQuestionRequest) (*dto.QuestionResponse, error) {
	// Verify quiz ownership
	quiz, err := s.repo.FindByID(ctx, quizID)
	if err != nil {
		return nil, err
	}

	if quiz.CreatedByID != userID {
		return nil, errors.New("you don't have permission to add questions to this quiz")
	}

	question := models.Question{
		QuizID:        quizID,
		QuestionText:  req.QuestionText,
		OptionA:       req.OptionA,
		OptionB:       req.OptionB,
		OptionC:       req.OptionC,
		OptionD:       req.OptionD,
		CorrectAnswer: req.CorrectAnswer,
		Points:        req.Points,
	}

	if err := s.repo.CreateQuestion(ctx, question); err != nil {
		return nil, err
	}

	return &dto.QuestionResponse{
		ID:           question.ID.String(),
		QuestionText: question.QuestionText,
		OptionA:      question.OptionA,
		OptionB:      question.OptionB,
		OptionC:      question.OptionC,
		OptionD:      question.OptionD,
		Points:       question.Points,
		OrderIndex:   question.OrderIndex,
	}, nil
}

// UpdateQuestion updates a question (owner only)
func (s *quizService) UpdateQuestion(ctx context.Context, questionID, userID uuid.UUID, req dto.CreateQuestionRequest) (*dto.QuestionResponse, error) {
	// TODO: Implement when FindQuestionByID is available in repository
	// For now, just update directly
	question := models.Question{
		QuestionText:  req.QuestionText,
		OptionA:       req.OptionA,
		OptionB:       req.OptionB,
		OptionC:       req.OptionC,
		OptionD:       req.OptionD,
		CorrectAnswer: req.CorrectAnswer,
		Points:        req.Points,
	}
	question.ID = questionID

	if err := s.repo.UpdateQuestion(ctx, question); err != nil {
		return nil, err
	}

	return &dto.QuestionResponse{
		ID:           question.ID.String(),
		QuestionText: question.QuestionText,
		OptionA:      question.OptionA,
		OptionB:      question.OptionB,
		OptionC:      question.OptionC,
		OptionD:      question.OptionD,
		Points:       question.Points,
		OrderIndex:   question.OrderIndex,
	}, nil
}

// DeleteQuestion deletes a question
func (s *quizService) DeleteQuestion(ctx context.Context, questionID, userID uuid.UUID) error {
	// TODO: Add ownership verification when FindQuestionByID is available
	return s.repo.DeleteQuestion(ctx, questionID)
}

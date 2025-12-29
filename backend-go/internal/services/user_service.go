package services

import (
	"context"
	"errors"
	"fmt"
	"quiz-game-backend/internal/dto"
	"quiz-game-backend/internal/models"
	"quiz-game-backend/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateProfile(ctx context.Context, id uuid.UUID, req dto.UpdateProfileRequest) (*models.User, error)
	ListUsers(ctx context.Context, page, limit int) (*dto.UserListResponse, error)
	UpdateStats(ctx context.Context, id uuid.UUID, req dto.UpdateStatsRequest) (*models.User, error)
	SetOnlineStatus(ctx context.Context, id uuid.UUID, isOnline bool) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (s *userService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.repo.FindByEmail(ctx, email)
}

func (s *userService) UpdateProfile(ctx context.Context, id uuid.UUID, req dto.UpdateProfileRequest) (*models.User, error) {
	// Get existing user
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	// Update fields
	if req.DisplayName != "" {
		user.DisplayName = req.DisplayName
	}
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}

	// Save
	if err := s.repo.UpdateProfile(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) ListUsers(ctx context.Context, page, limit int) (*dto.UserListResponse, error) {
	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit
	// Get users
	users, total, err := s.repo.ListUser(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	return &dto.UserListResponse{
		Users: dto.ToUserResponses(users),
		Total: int(total),
	}, nil
}

func (s *userService) UpdateStats(ctx context.Context, id uuid.UUID, req dto.UpdateStatsRequest) (*models.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	// Update stats
	user.ExperiencePoints += req.ExperiencePoints
	user.TotalGamesPlayed += req.GamesPlayed
	user.TotalWins += req.Wins
	// Calculate level (example: 100 XP per level)
	user.Level = (user.ExperiencePoints / 100) + 1
	if err := s.repo.UpdateProfile(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) SetOnlineStatus(ctx context.Context, id uuid.UUID, isOnline bool) error {
	return s.repo.UpdateOnlineStatus(ctx, id, isOnline)
}

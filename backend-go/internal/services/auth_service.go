package services

import (
	"context"
	"errors"
	"fmt"
	"quiz-game-backend/internal/config"
	"quiz-game-backend/internal/dto"
	"quiz-game-backend/internal/models"
	"quiz-game-backend/internal/repository"
	"quiz-game-backend/internal/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error)
	FacebookLogin(ctx context.Context, code string) (*dto.AuthResponse, error)
	GoogleLogin(ctx context.Context, code string) (*dto.AuthResponse, error)
	GetOAuthURL(provider string) string
}

type authService struct {
	userRepo     repository.UserRepository
	oauthService *OAuthService
	jwtConfig    *config.JWTConfig
}

func NewAuthService(
	userRepo repository.UserRepository,
	cfg *config.Config,
) AuthService { // Return interface
	return &authService{
		userRepo:     userRepo,
		oauthService: NewOAuthService(cfg.OAuth),
		jwtConfig:    cfg.JWT, // เก็บ config ทั้งหมด
	}
}

// Register
func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	// check email
	existingEmail, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingEmail != nil {
		return nil, fmt.Errorf("email already exists")
	}

	// check username
	existingUser, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, fmt.Errorf("username already exists")
	}

	// validate
	if err := utils.ValidatePasswordStrength(req.Password); err != nil {
		return nil, err
	}

	// hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// create user
	user := &models.User{
		Email:        req.Email,
		Username:     req.Username,
		DisplayName:  req.DisplayName,
		PasswordHash: hashedPassword,
		Level:        1,
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	//  jwt token
	token, err := utils.GenerateTokenRSA(
		user.ID,
		user.Email,
		s.jwtConfig.PrivateKey,
		s.jwtConfig.Expiration,
	)

	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		Token: token,
		User:  dto.ToUserResponse(user),
	}, nil

}

// login
func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	// find user by email
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("invalid email or password")
		}
		return nil, err
	}

	// check password
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, fmt.Errorf("invalid email or password")
	}

	now := time.Now()
	user.LastLoginAt = &now
	user.IsOnline = true
	s.userRepo.UpdateProfile(ctx, user)

	token, err := utils.GenerateTokenRSA(
		user.ID,
		user.Email,
		s.jwtConfig.PrivateKey,
		s.jwtConfig.Expiration,
	)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		Token: token,
		User:  dto.ToUserResponse(user),
	}, nil

}

func (s *authService) GetOAuthURL(provider string) string {
	state := uuid.New().String() // ควรเก็บไว้ใน session เพื่อ verify
	switch provider {
	case "facebook":
		return s.oauthService.GetAuthURL(ProviderFacebook, state)
	case "google":
		return s.oauthService.GetAuthURL(ProviderGoogle, state)
	}
	return ""
}

func (s *authService) FacebookLogin(ctx context.Context, code string) (*dto.AuthResponse, error) {
	// 1. Exchange code → token
	token, err := s.oauthService.ExchangeCode(ctx, ProviderFacebook, code)
	if err != nil {
		return nil, err
	}
	// 2. Get user info
	fbUser, err := s.oauthService.GetFacebookUserInfo(ctx, token)
	if err != nil {
		return nil, err
	}
	// 3. Find or create user
	user, _ := s.userRepo.FindByFacebookID(ctx, fbUser.ID)
	if user == nil {
		user = &models.User{
			Email:       fbUser.Email,
			Username:    "fb_" + fbUser.ID[:8],
			DisplayName: fbUser.Name,
			AvatarURL:   fbUser.Picture.Data.URL,
			FacebookID:  &fbUser.ID,
			Level:       1,
		}
		s.userRepo.CreateUser(ctx, user)
	}
	// 4. Generate JWT
	jwtToken, err := utils.GenerateTokenRSA(user.ID, user.Email, s.jwtConfig.PrivateKey, s.jwtConfig.Expiration)
	if err != nil {
		return nil, err
	}
	return &dto.AuthResponse{Token: jwtToken, User: dto.ToUserResponse(user)}, nil
}

func (s *authService) GoogleLogin(ctx context.Context, code string) (*dto.AuthResponse, error) {
	// 1. Exchange code → token
	token, err := s.oauthService.ExchangeCode(ctx, ProviderGoogle, code)
	if err != nil {
		return nil, err
	}
	// 2. Get user info
	googleUser, err := s.oauthService.GetGoogleUserInfo(ctx, token)
	if err != nil {
		return nil, err
	}
	// 3. Find or create user
	user, _ := s.userRepo.FindByGoogleID(ctx, googleUser.ID)
	if user == nil {
		user = &models.User{
			Email:       googleUser.Email,
			Username:    "google_" + googleUser.ID[:8],
			DisplayName: googleUser.Name,
			AvatarURL:   googleUser.Picture,
			GoogleId:    &googleUser.ID,
			Level:       1,
		}
		s.userRepo.CreateUser(ctx, user)
	}
	// 4. Generate JWT
	jwtToken, err := utils.GenerateTokenRSA(user.ID, user.Email, s.jwtConfig.PrivateKey, s.jwtConfig.Expiration)
	if err != nil {
		return nil, err
	}
	return &dto.AuthResponse{Token: jwtToken, User: dto.ToUserResponse(user)}, nil
}

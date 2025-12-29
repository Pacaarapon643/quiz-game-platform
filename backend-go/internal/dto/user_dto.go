package dto

import (
	"quiz-game-backend/internal/models"
	"time"
)

// ==================== Request DTOs ====================

// UpdateProfileRequest สำหรับแก้ไขโปรไฟล์
type UpdateProfileRequest struct {
	DisplayName string `json:"display_name" validate:"omitempty,min=1,max=50"`
	AvatarURL   string `json:"avatar_url" validate:"omitempty,url"`
}

// UpdateStatsRequest อัพเดทสถิติผู้เล่น
type UpdateStatsRequest struct {
	ExperiencePoints int `json:"experience_points" validate:"min=0"`
	GamesPlayed      int `json:"games_played" validate:"min=0"`
	Wins             int `json:"wins" validate:"min=0"`
}

// ==================== Response DTOs ====================
// UserResponse ข้อมูลผู้ใช้ (ไม่มี sensitive data)
type UserResponse struct {
	ID               string    `json:"id"`
	Email            string    `json:"email"`
	Username         string    `json:"username"`
	DisplayName      string    `json:"display_name"`
	AvatarURL        string    `json:"avatar_url"`
	Level            int       `json:"level"`
	ExperiencePoints int       `json:"experience_points"`
	TotalGamesPlayed int       `json:"total_games_played"`
	TotalWins        int       `json:"total_wins"`
	IsOnline         bool      `json:"is_online"`
	CreatedAt        time.Time `json:"created_at"`
}

// UserListResponse รายการผู้ใช้
type UserListResponse struct {
	Users []UserResponse `json:"users"`
	Total int            `json:"total"`
}

// ==================== Converters ====================
// ToUserResponse แปลง Model → DTO
func ToUserResponse(user *models.User) UserResponse {
	return UserResponse{
		ID:               user.ID.String(),
		Email:            user.Email,
		Username:         user.Username,
		DisplayName:      user.DisplayName,
		AvatarURL:        user.AvatarURL,
		Level:            user.Level,
		ExperiencePoints: user.ExperiencePoints,
		TotalGamesPlayed: user.TotalGamesPlayed,
		TotalWins:        user.TotalWins,
		IsOnline:         user.IsOnline,
		CreatedAt:        user.CreatedAt,
	}
}

// ToUserResponses แปลง []Model → []DTO
func ToUserResponses(users []models.User) []UserResponse {
	responses := make([]UserResponse, len(users))
	for i, user := range users {
		responses[i] = ToUserResponse(&user)
	}
	return responses
}

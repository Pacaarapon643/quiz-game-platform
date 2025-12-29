package models

import "time"

type User struct {
	BaseModel
	Email            string     `gorm:"uniqueIndex;not null" json:"email"`
	Username         string     `gorm:"uniqueIndex;not null" json:"username"`
	DisplayName      string     `gorm:"not null" json:"display_name"`
	PasswordHash     string     `gorm:"column:password_hash" json:"-"`
	AvatarURL        string     `json:"avatar_url"`
	FacebookID       *string    `gorm:"uniqueIndex" json:"facebook_id,omitempty"` // pointer = NULL support
	GoogleId         *string    `gorm:"uniqueIndex" json:"google_id,omitempty"`   // pointer = NULL support
	Level            int        `gorm:"default:1" json:"level"`
	ExperiencePoints int        `gorm:"default:0" json:"experience_points"`
	TotalGamesPlayed int        `gorm:"default:0" json:"total_games_played"`
	TotalWins        int        `gorm:"default:0" json:"total_wins"`
	IsOnline         bool       `gorm:"default:false" json:"is_online"`
	LastLoginAt      *time.Time `json:"last_login_at,omitempty"`

	// Relationships
	CreatedQuizzes []Quiz    `gorm:"foreignKey:CreatedByID" json:"created_quizzes,omitempty"`
	Messages       []Message `gorm:"foreignKey:UserID" json:"messages,omitempty"`
}

func (User) TableName() string {
	return "users"
}

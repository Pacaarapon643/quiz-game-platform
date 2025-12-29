package models

import "github.com/google/uuid"

type Quiz struct {
	BaseModel
	Title       string    `gorm:"not null" json:"title"`
	Description string    `json:"description"`
	Category    string    `gorm:"index" json:"category"`
	Difficulty  string    `gorm:"index" json:"difficulty"`      // easy, medium, hard
	TimeLimit   int       `gorm:"default:60" json:"time_limit"` // seconds
	IsPublished bool      `gorm:"default:false" json:"is_published"`
	CreatedByID uuid.UUID `gorm:"type:uuid;not null" json:"created_by_id"`

	// Relationships
	CreatedBy User       `gorm:"foreignKey:CreatedByID" json:"created_by,omitempty"`
	Questions []Question `gorm:"foreignKey:QuizID" json:"questions,omitempty"`
	GameRooms []GameRoom `gorm:"foreignKey:QuizID" json:"game_rooms,omitempty"`
}

func (Quiz) TableName() string {
	return "quizzes"
}

type Question struct {
	BaseModel
	QuizID        uuid.UUID `gorm:"type:uuid;not null;index" json:"quiz_id"`
	QuestionText  string    `gorm:"not null" json:"question_text"`
	OptionA       string    `gorm:"not null" json:"option_a"`
	OptionB       string    `gorm:"not null" json:"option_b"`
	OptionC       string    `gorm:"not null" json:"option_c"`
	OptionD       string    `gorm:"not null" json:"option_d"`
	CorrectAnswer string    `gorm:"not null" json:"-"` // A, B, C, D (ซ่อนไม่ส่งให้ client)
	Points        int       `gorm:"default:10" json:"points"`
	OrderIndex    int       `gorm:"default:0" json:"order_index"`

	// Relationships
	Quiz    Quiz         `gorm:"foreignKey:QuizID" json:"quiz,omitempty"`
	Answers []GameAnswer `gorm:"foreignKey:QuestionID" json:"answers,omitempty"`
}

func (Question) TableName() string {
	return "questions"
}

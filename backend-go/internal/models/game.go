package models

import (
	"time"

	"github.com/google/uuid"
)

type GameRoom struct {
	BaseModel
	RoomCode        string     `gorm:"uniqueIndex;not null" json:"room_code"`
	QuizID          uuid.UUID  `gorm:"type:uuid;not null;index" json:"quiz_id"`
	HostID          uuid.UUID  `gorm:"type:uuid;not null;index" json:"host_id"`
	MaxPlayers      int        `gorm:"default:10" json:"max_players"`
	Status          string     `gorm:"default:waiting" json:"status"` // waiting, playing, finished
	CurrentQuestion int        `gorm:"default:0" json:"current_question_index"`
	StartedAt       *time.Time `json:"started_at,omitempty"`
	FinishedAt      *time.Time `json:"finished_at,omitempty"`

	// Relationships
	Quiz         Quiz              `gorm:"foreignKey:QuizID" json:"quiz,omitempty"`
	Host         User              `gorm:"foreignKey:HostID" json:"host,omitempty"`
	Participants []GameParticipant `gorm:"foreignKey:GameRoomID" json:"participants,omitempty"`
	Messages     []Message         `gorm:"foreignKey:GameRoomID" json:"messages,omitempty"`
}

func (GameRoom) TableName() string {
	return "game_rooms"
}

type GameParticipant struct {
	BaseModel
	GameRoomID uuid.UUID `gorm:"type:uuid;not null;index" json:"game_room_id"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Score      int       `gorm:"default:0" json:"score"`
	Rank       int       `gorm:"default:0" json:"rank"`
	IsReady    bool      `gorm:"default:false" json:"is_ready"`
	JoinedAt   time.Time `json:"joined_at"`

	// Relationships
	GameRoom GameRoom     `gorm:"foreignKey:GameRoomID" json:"game_room,omitempty"`
	User     User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Answers  []GameAnswer `gorm:"foreignKey:ParticipantID" json:"answers,omitempty"`
}

func (GameParticipant) TableName() string {
	return "game_participants"
}

type GameAnswer struct {
	BaseModel
	ParticipantID uuid.UUID `gorm:"type:uuid;not null;index" json:"participant_id"`
	QuestionID    uuid.UUID `gorm:"type:uuid;not null;index" json:"question_id"`
	Answer        string    `gorm:"not null" json:"answer"` // A, B, C, D
	IsCorrect     bool      `gorm:"default:false" json:"is_correct"`
	TimeTaken     int       `json:"time_taken"` // milliseconds
	PointsEarned  int       `gorm:"default:0" json:"points_earned"`
	AnsweredAt    time.Time `json:"answered_at"`

	// Relationships
	Participant GameParticipant `gorm:"foreignKey:ParticipantID" json:"participant,omitempty"`
	Question    Question        `gorm:"foreignKey:QuestionID" json:"question,omitempty"`
}

func (GameAnswer) TableName() string {
	return "game_answers"
}

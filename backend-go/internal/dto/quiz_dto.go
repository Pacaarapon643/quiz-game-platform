package dto

// CreateQuizRequest
type CreateQuizRequest struct {
	Title       string `json:"title" validate:"required,min=3,max=100"`
	Description string `json:"description" validate:"max=500"`
	Category    string `json:"category" validate:"required"`
	Difficulty  string `json:"difficulty" validate:"required,oneof=easy medium hard"`
	TimeLimit   int    `json:"time_limit" validate:"min=10,max=300"`
}

// UpdateQuizRequest
type UpdateQuizRequest struct {
	Title       string `json:"title" validate:"omitempty,min=3,max=100"`
	Description string `json:"description" validate:"max=500"`
	Category    string `json:"category"`
	Difficulty  string `json:"difficulty" validate:"omitempty,oneof=easy medium hard"`
	TimeLimit   int    `json:"time_limit" validate:"omitempty,min=10,max=300"`
	IsPublished bool   `json:"is_published"`
}

// CreateQuestionRequest
type CreateQuestionRequest struct {
	QuestionText  string `json:"question_text" validate:"required"`
	OptionA       string `json:"option_a" validate:"required"`
	OptionB       string `json:"option_b" validate:"required"`
	OptionC       string `json:"option_c" validate:"required"`
	OptionD       string `json:"option_d" validate:"required"`
	CorrectAnswer string `json:"correct_answer" validate:"required,oneof=A B C D"`
	Points        int    `json:"points" validate:"min=1,max=100"`
}

// QuizResponse
type QuizResponse struct {
	ID          string             `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Category    string             `json:"category"`
	Difficulty  string             `json:"difficulty"`
	TimeLimit   int                `json:"time_limit"`
	IsPublished bool               `json:"is_published"`
	CreatedByID string             `json:"created_by_id"`
	Questions   []QuestionResponse `json:"questions,omitempty"`
}

// QuestionResponse
type QuestionResponse struct {
	ID           string `json:"id"`
	QuestionText string `json:"question_text"`
	OptionA      string `json:"option_a"`
	OptionB      string `json:"option_b"`
	OptionC      string `json:"option_c"`
	OptionD      string `json:"option_d"`
	Points       int    `json:"points"`
	OrderIndex   int    `json:"order_index"`
}

// QuizListResponse for paginated quiz list
type QuizListResponse struct {
	Quizzes []QuizResponse `json:"quizzes"`
	Total   int64          `json:"total"`
}

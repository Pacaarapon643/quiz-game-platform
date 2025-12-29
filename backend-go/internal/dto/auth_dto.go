package dto

// RegisterRequest สำหรับสมัครสมาชิก
type RegisterRequest struct {
	Email       string `json:"email" validate:"required,email,max=100"`
	Username    string `json:"username" validate:"required,min=3,max=20,alphanum"`
	DisplayName string `json:"display_name" validate:"required,min=1,max=50"`
	Password    string `json:"password" validate:"required,min=8,max=100"`
}

// LoginRequest สำหรับเข้าสู่ระบบ
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// FacebookCallbackRequest จาก Facebook OAuth
type FacebookCallbackRequest struct {
	Code string `json:"code" validate:"required"`
}

// AuthResponse response หลัง login/register
type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

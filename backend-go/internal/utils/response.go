package utils

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// Response โครงสร้าง response มาตรฐาน
type Response struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
	Timestamp string      `json:"timestamp"`
}

// PaginatedResponse โครงสร้างสำหรับข้อมูลแบบ pagination
type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
	RequestID  string      `json:"request_id,omitempty"`
	Timestamp  string      `json:"timestamp"`
}

// Pagination ข้อมูล pagination
type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalPages int `json:"total_pages"`
	TotalItems int `json:"total_items"`
}

// ValidationErrorItem รายละเอียด validation error
type ValidationErrorItem struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Success ส่ง response สำเร็จพร้อมข้อมูล
func Success(c *fiber.Ctx, data interface{}, message ...string) error {
	msg := "Success"
	if len(message) > 0 {
		msg = message[0]
	}

	return c.JSON(Response{
		Success:   true,
		Message:   msg,
		Data:      data,
		RequestID: getRequestID(c),
		Timestamp: time.Now().Format(time.RFC3339),
	})
}

// SuccessWithStatus ส่ง response สำเร็จพร้อม custom status code
func SuccessWithStatus(c *fiber.Ctx, status int, data interface{}, message ...string) error {
	msg := "Success"
	if len(message) > 0 {
		msg = message[0]
	}

	return c.Status(status).JSON(Response{
		Success:   true,
		Message:   msg,
		Data:      data,
		RequestID: getRequestID(c),
		Timestamp: time.Now().Format(time.RFC3339),
	})
}

// Created ส่ง response เมื่อสร้างข้อมูลสำเร็จ (201)
func Created(c *fiber.Ctx, data interface{}, message ...string) error {
	msg := "Resource created successfully"
	if len(message) > 0 {
		msg = message[0]
	}

	return c.Status(fiber.StatusCreated).JSON(Response{
		Success:   true,
		Message:   msg,
		Data:      data,
		RequestID: getRequestID(c),
		Timestamp: time.Now().Format(time.RFC3339),
	})
}

// NoContent ส่ง response 204 (ไม่มีข้อมูล)
func NoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

// Error ส่ง response error ทั่วไป
func Error(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(Response{
		Success:   false,
		Error:     message,
		RequestID: getRequestID(c),
		Timestamp: time.Now().Format(time.RFC3339),
	})
}

// BadRequest ส่ง response 400
func BadRequest(c *fiber.Ctx, message ...string) error {
	msg := "Bad request"
	if len(message) > 0 {
		msg = message[0]
	}

	return Error(c, fiber.StatusBadRequest, msg)
}

// Unauthorized ส่ง response 401
func Unauthorized(c *fiber.Ctx, message ...string) error {
	msg := "Unauthorized"
	if len(message) > 0 {
		msg = message[0]
	}

	return Error(c, fiber.StatusUnauthorized, msg)
}

// Forbidden ส่ง response 403
func Forbidden(c *fiber.Ctx, message ...string) error {
	msg := "Forbidden"
	if len(message) > 0 {
		msg = message[0]
	}

	return Error(c, fiber.StatusForbidden, msg)
}

// NotFound ส่ง response 404
func NotFound(c *fiber.Ctx, message ...string) error {
	msg := "Resource not found"
	if len(message) > 0 {
		msg = message[0]
	}

	return Error(c, fiber.StatusNotFound, msg)
}

// Conflict ส่ง response 409
func Conflict(c *fiber.Ctx, message ...string) error {
	msg := "Resource already exists"
	if len(message) > 0 {
		msg = message[0]
	}

	return Error(c, fiber.StatusConflict, msg)
}

// InternalServerError ส่ง response 500
func InternalServerError(c *fiber.Ctx, message ...string) error {
	msg := "Internal server error"
	if len(message) > 0 {
		msg = message[0]
	}

	return Error(c, fiber.StatusInternalServerError, msg)
}

// ValidationError ส่ง response สำหรับ validation errors
func ValidationError(c *fiber.Ctx, errors []ValidationErrorItem) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"success":    false,
		"error":      "Validation failed",
		"details":    errors,
		"request_id": getRequestID(c),
		"timestamp":  time.Now().Format(time.RFC3339),
	})
}

// Paginated ส่ง response พร้อม pagination
func Paginated(c *fiber.Ctx, data interface{}, page, limit, totalItems int) error {
	totalPages := (totalItems + limit - 1) / limit

	return c.JSON(PaginatedResponse{
		Success: true,
		Data:    data,
		Pagination: Pagination{
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
			TotalItems: totalItems,
		},
		RequestID: getRequestID(c),
		Timestamp: time.Now().Format(time.RFC3339),
	})
}

// getRequestID ดึง request ID จาก context
func getRequestID(c *fiber.Ctx) string {
	if requestID := c.Locals("requestid"); requestID != nil {
		return requestID.(string)
	}
	return ""
}

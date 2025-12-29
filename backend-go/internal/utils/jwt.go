package utils

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateToken สร้าง JWT token ด้วย RSA private key
func GenerateTokenRSA(userID uuid.UUID, email string, privateKey *rsa.PrivateKey, expiration time.Duration) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID: userID.String(),
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        uuid.New().String(),
		},
	}
	// ใช้ RS256 algorithm
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}

// ValidateTokenRSA ตรวจสอบ JWT token ด้วย RSA public key
func ValidateTokenRSA(tokenString string, publicKey *rsa.PublicKey) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// ตรวจสอบว่าใช้ RSA algorithm
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}

// ExtractUserIDRSA ดึง User ID จาก token
func ExtractUserIDRSA(tokenString string, publicKey *rsa.PublicKey) (uuid.UUID, error) {
	claims, err := ValidateTokenRSA(tokenString, publicKey)
	if err != nil {
		return uuid.Nil, err
	}
	return uuid.Parse(claims.UserID)
}

// SetAuthCookie ตั้งค่า JWT token ใน HTTP-Only cookie
func SetAuthCookie(c *fiber.Ctx, token string, expiration time.Duration, isProduction bool) {
	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    token,
		Expires:  time.Now().Add(expiration),
		HTTPOnly: true,
		Secure:   false, // ปิด secure สำหรับ localhost (HTTP)
		SameSite: "Lax",
		Path:     "/",
	})
}

// ClearAuthCookie ลบ auth cookie
func ClearAuthCookie(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Path:     "/",
	})
}

// GetTokenFromCookie ดึง token จาก cookie
func GetTokenFromCookie(c *fiber.Ctx) string {
	return c.Cookies("auth_token")
}

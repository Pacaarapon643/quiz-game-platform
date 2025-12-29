package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type JWTConfig struct {
	PrivateKeyPath string        `validate:"required"`
	PublicKeyPath  string        `validate:"required"`
	Expiration     time.Duration `validate:"required"`

	// Loaded keys
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

func LoadJWTConfig() *JWTConfig {
	expiration := viper.GetDuration("JWT_EXPIRATION")

	// ถ้า expiration เป็น 0 ให้ใช้ค่า default 24 ชั่วโมง
	if expiration == 0 {
		expiration = 24 * time.Hour
		log.Println("⚠️ JWT_EXPIRATION not set, using default: 24h")
	}
	log.Printf("📌 JWT Expiration: %v", expiration)

	jwtConfig := JWTConfig{
		PrivateKeyPath: viper.GetString("JWT_PRIVATE_KEY_PATH"),
		PublicKeyPath:  viper.GetString("JWT_PUBLIC_KEY_PATH"),
		Expiration:     expiration,
	}

	validate := validator.New()
	if err := validate.Struct(jwtConfig); err != nil {
		log.Fatal("Invalid config: ", err)
	}

	// โหลด private key
	privateKey, err := loadPrivateKey(jwtConfig.PrivateKeyPath)
	if err != nil {
		log.Fatal("Failed to load private key: ", err)
	}
	jwtConfig.PrivateKey = privateKey

	// โหลด public key
	publicKey, err := loadPublicKey(jwtConfig.PublicKeyPath)
	if err != nil {
		log.Fatal("Failed to load public key: ", err)
	}
	jwtConfig.PublicKey = publicKey
	return &jwtConfig

}

func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	keyData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}
	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}
	return privateKey, nil
}

func loadPublicKey(path string) (*rsa.PublicKey, error) {
	keyData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key: %w", err)
	}
	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}
	publicKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA public key")
	}
	return publicKey, nil
}

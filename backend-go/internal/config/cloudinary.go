package config

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type CloudinaryConfig struct {
	CloudName string `validate:"required"`
	APIKey    string `validate:"required"`
	APISecret string `validate:"required"`
}

func LoadCloudinaryConfig() *CloudinaryConfig {
	cloudinaryConfig := CloudinaryConfig{
		CloudName: viper.GetString("CLOUDINARY_CLOUD_NAME"),
		APIKey:    viper.GetString("CLOUDINARY_API_KEY"),
		APISecret: viper.GetString("CLOUDINARY_API_SECRET"),
	}

	validate := validator.New()
	if err := validate.Struct(cloudinaryConfig); err != nil {
		log.Fatal("Invalid config: ", err)
	}

	return &cloudinaryConfig

}

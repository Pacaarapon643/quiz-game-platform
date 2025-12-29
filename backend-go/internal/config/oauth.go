package config

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type OAuthConfig struct {
	// Facebook
	FacebookAppID       string `validate:"required"`
	FacebookAppSecret   string `validate:"required"`
	FacebookRedirectURL string `validate:"required"`

	//Google
	GoogleClientID     string `validate:"required"`
	GoogleClientSecret string `validate:"required"`
	GoogleRedirectURL  string `validate:"required"`
}

func LoadOAuthConfig() *OAuthConfig {
	oAuthConfig := OAuthConfig{
		FacebookAppID:       viper.GetString("FACEBOOK_APP_ID"),
		FacebookAppSecret:   viper.GetString("FACEBOOK_APP_SECRET"),
		FacebookRedirectURL: viper.GetString("FACEBOOK_REDIRECT_URL"),

		GoogleClientID:     viper.GetString("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: viper.GetString("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL:  viper.GetString("GOOGLE_REDIRECT_URL"),
	}

	validator := validator.New()
	if err := validator.Struct(oAuthConfig); err != nil {
		log.Fatal("Invalid config: ", err)
	}

	return &oAuthConfig
}

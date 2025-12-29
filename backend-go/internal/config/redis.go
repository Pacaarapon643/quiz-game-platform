package config

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type RedisConfig struct {
	URL string `validate:"required"`
	DB  int    `validate:"min=0"`
}

func LoadRedisConfig() *RedisConfig {
	redisConfig := RedisConfig{
		URL: viper.GetString("REDIS_URL"),
		DB:  viper.GetInt("REDIS_DB"),
	}

	validate := validator.New()
	if err := validate.Struct(redisConfig); err != nil {
		log.Fatal("Invalid config: ", err)
	}

	return &redisConfig
}

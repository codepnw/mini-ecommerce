package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/codepnw/mini-ecommerce/pkg/validate"
	"github.com/joho/godotenv"
)

type EnvConfig struct {
	APP AppConfig `envPrefix:"APP_"`
	DB  DBConfig  `envPrefix:"DB_"`
	JWT JWTConfig `envPrefix:"JWT_"`
}

type AppConfig struct {
	Version int `env:"VERSION" envDefault:"1"`
	Port    int `env:"PORT" envDefault:"8080"`
}

type DBConfig struct {
	User     string `env:"USER" validate:"required"`
	Password string `env:"PASSWORD"`
	Host     string `env:"HOST" envDefault:"127.0.0.1"`
	Port     int    `env:"PORT" envDefault:"5432"`
	DBName   string `env:"NAME" validate:"required"`
	SSLMode  string `env:"SSL_MODE" envDefault:"disable"`
}

type JWTConfig struct {
	SecretKey  string `env:"SECRET_KEY" validate:"required"`
	RefreshKey string `env:"REFRESH_KEY" validate:"required"`
}

func LoadConfig(path string) (*EnvConfig, error) {
	if err := godotenv.Load(path); err != nil {
		/*
			NOTE: No Docker return error
			return nil, fmt.Errorf("load env failed: %w", err)

			NOTE: Docker copy .env variables to OS environment variables.
		*/
		log.Println("⚠️  Warning: .env file not found. Using OS environment variables.")
	}

	cfg := new(EnvConfig)
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parse env failed: %w", err)
	}

	if err := validate.Struct(cfg); err != nil {
		return nil, fmt.Errorf("validate env failed: %w", err)
	}
	return cfg, nil
}

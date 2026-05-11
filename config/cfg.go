package config

import (
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

type Config struct {
	googleAuth oauth2.Config
}

func NewConfig() (*Config, error) {
	if os.Getenv("ENVIRONMENT") != "DEV" {
		if err := godotenv.Load(); err != nil {
			return &Config{}, err
		}
	}

	googleOAuth2Cfg, err := FillOAuthGoogleCfg()

	return &Config{googleAuth: googleOAuth2Cfg}, err
}

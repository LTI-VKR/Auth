package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GoogleAuth     OAuth2Cfg
	OAuthSecretKey string
}

func NewConfig() (*Config, error) {
	if os.Getenv("ENVIRONMENT") != "DEV" {
		if err := godotenv.Load(); err != nil {
			return &Config{}, err
		}
	}

	googleOAuth2Cfg, err := GoogleCfg()
	if err != nil {
		return &Config{}, err
	}
	oauthSecretKey := os.Getenv("OAUTH_SECRET_KEY")

	return &Config{GoogleAuth: googleOAuth2Cfg, OAuthSecretKey: oauthSecretKey}, err
}

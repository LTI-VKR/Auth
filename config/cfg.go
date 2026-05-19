package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GoogleAuth     OAuth2Cfg
	VkAuth         OAuth2Cfg
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
		panic("Error loading Google OAuth2 config")
	}
	vkOAuth2Cfg, err := VKCfg()
	if err != nil {
		panic("Error loading VK OAuth2 config")
	}
	oauthSecretKey := os.Getenv("OAUTH_SECRET_KEY")

	return &Config{
		GoogleAuth:     googleOAuth2Cfg,
		VkAuth:         vkOAuth2Cfg,
		OAuthSecretKey: oauthSecretKey,
	}, err
}

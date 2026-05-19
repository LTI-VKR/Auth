package config

import (
	"os"
)

type OAuth2Cfg struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	AuthURL      string
	TokenURL     string
	InfoURL      string
	Scopes       []string
}

func GoogleCfg() (OAuth2Cfg, error) {
	var googleOauthConfig = OAuth2Cfg{
		ClientID:     os.Getenv("OAUTH_GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH_GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("OAUTH_GOOGLE_CALLBACK_URL"),
		AuthURL:      os.Getenv("OAUTH_GOOGLE_AUTH_URL"),
		TokenURL:     os.Getenv("OAUTH_GOOGLE_TOKEN_URL"),
		Scopes:       []string{"openid", "email", "profile"},
	}

	return googleOauthConfig, nil
}

func VKCfg() (OAuth2Cfg, error) {
	var vkOauthConfig = OAuth2Cfg{
		ClientID:     os.Getenv("OAUTH_VK_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH_VK_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("OAUTH_VK_CALLBACK_URL"),
		AuthURL:      os.Getenv("OAUTH_VK_AUTH_URL"),
		TokenURL:     os.Getenv("OAUTH_VK_TOKEN_URL"),
		InfoURL:      os.Getenv("OAUTH_VK_INFO_URL"),
		Scopes:       []string{"vkid.personal_info", "email"},
	}

	return vkOauthConfig, nil
}

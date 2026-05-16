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
	Scopes       []string
}

func GoogleCfg() (OAuth2Cfg, error) {
	var googleOauthConfig = OAuth2Cfg{
		ClientID:     os.Getenv("OAUTH_GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH_GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("OAUTH_GOOGLE_CALLBACK_URL"),
		Scopes:       []string{"openid", "email", "profile"},
	}

	return googleOauthConfig, nil
}

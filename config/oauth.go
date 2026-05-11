package config

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func FillOAuthGoogleCfg() (oauth2.Config, error) {
	var googleOauthConfig = oauth2.Config{
		RedirectURL:  os.Getenv("OAUTH_GOOGLE_CALLBACK_URL"),
		ClientID:     os.Getenv("OAUTH_GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH_GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     google.Endpoint,
	}

	return googleOauthConfig, nil
}

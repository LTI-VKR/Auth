package entity

import "time"

type AuthToken struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

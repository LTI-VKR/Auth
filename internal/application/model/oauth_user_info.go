package model

import "auth/internal/application"

type OAuthUserInfo struct {
	Provider      application.OAuthProvider
	Id            string
	Email         string
	Name          string
	FirstName     string
	LastName      string
	AvatarURL     string
	EmailVerified bool
}

func NewGoogleUserInfo(email, name, givenName, familyName, SubId, picture string, emailVerified bool, provider application.OAuthProvider) OAuthUserInfo {
	return OAuthUserInfo{
		Provider:      provider,
		Id:            SubId,
		Email:         email,
		Name:          name,
		FirstName:     givenName,
		LastName:      familyName,
		AvatarURL:     picture,
		EmailVerified: emailVerified,
	}
}

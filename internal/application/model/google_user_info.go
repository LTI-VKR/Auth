package model

type GoogleUserInfo struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	SubId         string `json:"sub"`
	Picture       string `json:"picture"`
}

func NewGoogleUserInfo(email, name, givenName, familyName, SubId, picture string, emailVerified bool) GoogleUserInfo {
	return GoogleUserInfo{
		Email:         email,
		EmailVerified: emailVerified,
		Name:          name,
		GivenName:     givenName,
		FamilyName:    familyName,
		SubId:         SubId,
		Picture:       picture,
	}
}

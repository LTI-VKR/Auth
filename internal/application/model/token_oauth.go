package model

type TokenOAuth struct {
	AccessToken  string //
	RefreshToken string //
	IdToken      string //
	UserId       string
	State        string
	Scope        string //
	TokenType    string //
	ExpiresIn    int    //
	ClientId     string
}

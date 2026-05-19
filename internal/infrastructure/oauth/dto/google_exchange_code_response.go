package dto

import "auth/internal/application/model"

type GoogleExchangeCodeResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	IdToken      string `json:"id_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func (gr *GoogleExchangeCodeResponse) ToTokenOAuthModel(clientId string) model.TokenOAuth {
	return model.TokenOAuth{
		AccessToken:  gr.AccessToken,
		RefreshToken: gr.RefreshToken,
		IdToken:      gr.IdToken,
		UserId:       "",
		State:        "",
		Scope:        gr.Scope,
		TokenType:    gr.TokenType,
		ExpiresIn:    gr.ExpiresIn,
		ClientId:     clientId,
	}
}

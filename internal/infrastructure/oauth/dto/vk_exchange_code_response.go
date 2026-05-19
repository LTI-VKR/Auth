package dto

import (
	"auth/internal/application/model"
	"encoding/json"
)

type VkExchangeCodeResponse struct {
	RefreshToken string      `json:"refresh_token"`
	AccessToken  string      `json:"access_token"`
	IdToken      string      `json:"id_token"`
	TokenType    string      `json:"token_type"`
	ExpiresIn    int         `json:"expires_in"`
	UserId       json.Number `json:"user_id"`
	State        string      `json:"state"`
	Scope        string      `json:"scope"`
}

func (vr *VkExchangeCodeResponse) ToTokenOAuthModel(clientId string) model.TokenOAuth {
	return model.TokenOAuth{
		AccessToken:  vr.AccessToken,
		RefreshToken: vr.RefreshToken,
		IdToken:      vr.IdToken,
		UserId:       string(vr.UserId),
		State:        vr.State,
		Scope:        vr.Scope,
		TokenType:    vr.TokenType,
		ExpiresIn:    vr.ExpiresIn,
		ClientId:     clientId,
	}
}

package query

import (
	"auth/internal/application/model"
	"auth/internal/application/ports"
	"context"
)

type GetGoogleOAuthCallbackQuery struct {
	oauthClient ports.OAuthClient
}

func NewGetGoogleOAuthCallbackQuery(client ports.OAuthClient) *GetGoogleOAuthCallbackQuery {
	return &GetGoogleOAuthCallbackQuery{
		oauthClient: client,
	}
}

func (c *GetGoogleOAuthCallbackQuery) Execute(ctx context.Context, code string) (model.GoogleUserInfo, error) {
	userInfo, err := c.oauthClient.ExchangeCode(ctx, code)
	if err != nil {
		return model.GoogleUserInfo{}, err
	}

	return userInfo, nil
}

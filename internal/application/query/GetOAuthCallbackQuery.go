package query

import (
	"auth/internal/application"
	"auth/internal/application/model"
	"auth/internal/application/ports"
	"context"
)

type GetOAuthCallbackQuery struct {
	clientFactory ports.OAuthClientFactory
}

func NewGetOAuthCallbackQuery(factory ports.OAuthClientFactory) *GetOAuthCallbackQuery {
	return &GetOAuthCallbackQuery{clientFactory: factory}
}

func (c *GetOAuthCallbackQuery) Execute(ctx context.Context, params model.ExchangeCodeParams, provider application.OAuthProvider) (model.OAuthUserInfo, error) {
	client, err := c.clientFactory.Get(provider)
	if err != nil {
		return model.OAuthUserInfo{}, err
	}

	token, err := client.ExchangeCode(ctx, params)
	oauthUserData, err := client.GetUserInfo(ctx, token) // здесь конечные данные из oauth

	return oauthUserData, nil
}

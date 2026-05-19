package query

import (
	"auth/internal/application"
	"auth/internal/application/ports"
	"context"
)

type GetOAuthRedirectUrlQuery struct {
	clientFactory ports.OAuthClientFactory
}

func NewGetOAuthRedirectUrlQuery(factory ports.OAuthClientFactory) *GetOAuthRedirectUrlQuery {
	return &GetOAuthRedirectUrlQuery{clientFactory: factory}
}

func (c *GetOAuthRedirectUrlQuery) Execute(ctx context.Context, state string, oauthProvider application.OAuthProvider) (string, error) {
	client, err := c.clientFactory.Get(oauthProvider)
	if err != nil {
		return "", err
	}

	return client.AuthCodeURL(ctx, state)
}

package query

import (
	"auth/internal/application/ports"
	"context"
)

type GetGoogleOAuthRedirectQuery struct {
	oauthClient ports.OAuthClient
}

func NewGetGoogleOAuthRedirectQuery(client ports.OAuthClient) *GetGoogleOAuthRedirectQuery {
	return &GetGoogleOAuthRedirectQuery{
		oauthClient: client,
	}
}

func (c *GetGoogleOAuthRedirectQuery) Execute(ctx context.Context, state string) (string, error) {
	URL, err := c.oauthClient.AuthCodeURL(ctx, state)
	if err != nil {
		return "", err
	}

	return URL, nil
}

package ports

import (
	"auth/internal/application/model"
	"context"
)

type OAuthClient interface {
	AuthCodeURL(ctx context.Context, state string) (string, error)
	ExchangeCode(ctx context.Context, params model.ExchangeCodeParams) (model.TokenOAuth, error)
	GetUserInfo(ctx context.Context, token model.TokenOAuth) (model.OAuthUserInfo, error)
}

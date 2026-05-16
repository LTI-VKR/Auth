package ports

import (
	"auth/internal/application/model"
	"context"
)

type OAuthClient interface {
	AuthCodeURL(ctx context.Context, state string) (string, error)
	ExchangeCode(ctx context.Context, code string) (model.GoogleUserInfo, error)
}

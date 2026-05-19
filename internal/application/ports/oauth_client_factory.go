package ports

import "auth/internal/application"

type OAuthClientFactory interface {
	Get(provider application.OAuthProvider) (OAuthClient, error)
}

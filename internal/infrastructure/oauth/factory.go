package oauth

import (
	"auth/internal/application"
	"auth/internal/application/ports"
)

type ClientFactory struct {
	clients map[application.OAuthProvider]ports.OAuthClient
}

func NewClientFactory(clients map[application.OAuthProvider]ports.OAuthClient) *ClientFactory {
	return &ClientFactory{clients: clients}
}

func (f *ClientFactory) Get(provider application.OAuthProvider) (ports.OAuthClient, error) {
	client, ok := f.clients[provider]
	if !ok {
		return nil, ErrInvalidOAuthProvider
	}

	return client, nil
}

package oauth

import (
	"auth/config"
	"auth/internal/application"
	"auth/internal/application/model"
	"auth/internal/application/ports"
	"auth/internal/infrastructure/oauth/dto"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2"
)

type vkClient struct {
	config     *oauth2.Config
	httpClient *http.Client
}

func NewVkClient(cfg *config.OAuth2Cfg, client *http.Client) ports.OAuthClient {
	return &vkClient{
		config: &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.RedirectURL,
			Scopes:       cfg.Scopes,
			Endpoint: oauth2.Endpoint{
				AuthURL:       cfg.AuthURL,
				TokenURL:      cfg.TokenURL,
				DeviceAuthURL: cfg.InfoURL,
			},
		},
		httpClient: client,
	}
}

// AuthCodeURL составляет и возвращает ссылку на OAuth
func (c *vkClient) AuthCodeURL(ctx context.Context, state string) (string, error) {
	_ = ctx
	challenge := oauth2.S256ChallengeOption(state)
	authURL := c.config.AuthCodeURL(state, oauth2.AccessTypeOffline, challenge, oauth2.SetAuthURLParam("lang_id", "0"))
	return authURL, nil
}

// ExchangeCode обменивает authorization code на токены
func (c *vkClient) ExchangeCode(ctx context.Context, params model.ExchangeCodeParams) (model.TokenOAuth, error) {
	codeVerifier := params.CodeVerifier
	if codeVerifier == "" {
		codeVerifier = params.State
	}

	val := url.Values{
		"grant_type":    {"authorization_code"},
		"code_verifier": {codeVerifier},
		"redirect_uri":  {c.config.RedirectURL},
		"code":          {params.Code},
		"client_id":     {c.config.ClientID},
		"device_id":     {params.DeviceID},
		"state":         {params.State},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.config.Endpoint.TokenURL, strings.NewReader(val.Encode()))
	if err != nil {
		return model.TokenOAuth{}, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return model.TokenOAuth{}, err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return model.TokenOAuth{}, fmt.Errorf("ошибка при обмене кода на токены: %s", body)
	}

	var tokenResp dto.VkExchangeCodeResponse
	if err = json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return model.TokenOAuth{}, fmt.Errorf("неверный ответ от сервера VK: %v", err)
	}

	//TODO сделать проверку на все данные

	return tokenResp.ToTokenOAuthModel(c.config.ClientID), nil
}

// GetUserInfo Отдельный запрос для получения данных пользователя
func (c *vkClient) GetUserInfo(ctx context.Context, token model.TokenOAuth) (model.OAuthUserInfo, error) {
	val := url.Values{
		"access_token": {token.AccessToken},
		"client_id":    {token.ClientId},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.config.Endpoint.DeviceAuthURL, strings.NewReader(val.Encode()))
	if err != nil {
		return model.OAuthUserInfo{}, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return model.OAuthUserInfo{}, err
	}
	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return model.OAuthUserInfo{}, fmt.Errorf("ошибка при получении данных пользователя: %s", body)
	}

	var userinfoResp dto.VkUserinfoResponse
	if err = json.NewDecoder(resp.Body).Decode(&userinfoResp); err != nil {
		return model.OAuthUserInfo{}, err
	}

	return model.OAuthUserInfo{
		Provider:      application.VkProvider,
		Id:            userinfoResp.User.UserId,
		Email:         userinfoResp.User.Email,
		Name:          userinfoResp.User.FirstName + " " + userinfoResp.User.LastName,
		FirstName:     userinfoResp.User.FirstName,
		LastName:      userinfoResp.User.LastName,
		AvatarURL:     userinfoResp.User.Avatar,
		EmailVerified: true,
	}, nil
}

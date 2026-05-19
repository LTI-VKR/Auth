package oauth

import (
	"auth/config"
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

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

type googleClient struct {
	config     *oauth2.Config
	httpClient *http.Client
}

func NewGoogleClient(cfg *config.OAuth2Cfg, client *http.Client) ports.OAuthClient {
	return &googleClient{
		config: &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.RedirectURL,
			Scopes:       cfg.Scopes,
			Endpoint: oauth2.Endpoint{
				AuthURL:  cfg.AuthURL,
				TokenURL: cfg.TokenURL,
			},
		},
		httpClient: client,
	}
}

// AuthCodeURL составляет и возвращает ссылку на OAuth
func (c *googleClient) AuthCodeURL(ctx context.Context, state string) (string, error) {
	_ = ctx
	authURL := c.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
	return authURL, nil
}

// ExchangeCode обменивает authorization code на токены
func (c *googleClient) ExchangeCode(ctx context.Context, params model.ExchangeCodeParams) (model.TokenOAuth, error) {
	val := url.Values{
		"client_id":     {c.config.ClientID},
		"client_secret": {c.config.ClientSecret},
		"code":          {params.Code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {c.config.RedirectURL},
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

	var tokenResp dto.GoogleExchangeCodeResponse
	if err = json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return model.TokenOAuth{}, fmt.Errorf("неверный ответ от сервера Google: %v", err)
	}
	//TODO сделать проверку на все данные

	return tokenResp.ToTokenOAuthModel(c.config.ClientID), nil
}

// GetUserInfo Отдельный запрос для получения данных пользователя
func (c *googleClient) GetUserInfo(ctx context.Context, tokenOAuth model.TokenOAuth) (model.OAuthUserInfo, error) {
	_ = ctx
	token, _, err := new(jwt.Parser).ParseUnverified(tokenOAuth.IdToken, jwt.MapClaims{})
	if err != nil {
		return model.OAuthUserInfo{}, fmt.Errorf("%w: %v", ErrParseIDToken, err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return model.OAuthUserInfo{}, fmt.Errorf("неверный формат ID Token claims")
	}

	email, _ := claims["email"].(string)
	emailVerified, _ := claims["email_verified"].(bool)
	name, _ := claims["name"].(string)
	sub, _ := claims["sub"].(string)
	familyName, _ := claims["family_name"].(string)
	givenName, _ := claims["given_name"].(string)
	picture, _ := claims["picture"].(string)

	// Проверка обязательных полей
	var missingFields []string
	if email == "" {
		missingFields = append(missingFields, "email")
	}
	if name == "" {
		missingFields = append(missingFields, "name")
	}
	if sub == "" {
		missingFields = append(missingFields, "sub")
	}

	if len(missingFields) > 0 {
		errMsg := fmt.Sprintf("отсутствуют поля: %s", strings.Join(missingFields, ", "))
		return model.OAuthUserInfo{}, fmt.Errorf("отсутствуют обязательные поля в ID Token: %s", errMsg)
	}

	userInfo := model.NewGoogleUserInfo(email, name, givenName, familyName, sub, picture, emailVerified, "")
	return userInfo, nil
}

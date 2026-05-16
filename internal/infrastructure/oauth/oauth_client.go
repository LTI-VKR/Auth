package oauth

import (
	"auth/config"
	"auth/internal/application/model"
	"auth/internal/infrastructure/oauth/dto"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

const ClientTimeout = 30 * time.Second

// Ошибки OAuth

type OAuthClient struct {
	config     *oauth2.Config
	httpClient *http.Client
}

func NewOAuth2Client(cfg *config.OAuth2Cfg) *OAuthClient {
	return &OAuthClient{
		config: &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.RedirectURL,
			Scopes:       []string{"openid", "profile", "email"},
			Endpoint:     endpoints.Google,
		},
		httpClient: &http.Client{Timeout: ClientTimeout},
	}
}

// AuthCodeURL составляет и возвращает ссылку на Google OAuth
func (c *OAuthClient) AuthCodeURL(ctx context.Context, state string) (string, error) {
	authURL := c.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
	return authURL, nil
}

// ExchangeCode обменивает authorization code на токены
func (c *OAuthClient) ExchangeCode(ctx context.Context, code string) (model.GoogleUserInfo, error) {
	val := url.Values{
		"client_id":     {c.config.ClientID},
		"client_secret": {c.config.ClientSecret},
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {c.config.RedirectURL},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.config.Endpoint.TokenURL, strings.NewReader(val.Encode()))
	if err != nil {
		log.Printf("failed to create request: %v", err)
		return model.GoogleUserInfo{}, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Printf("failed to send request: %v", err)
		return model.GoogleUserInfo{}, err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}(resp.Body)

	// Проверка статуса ответа
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("oauth server error: status=%d body=%s", resp.StatusCode, body)
		return model.GoogleUserInfo{}, fmt.Errorf("ошибка при обмене кода на токены: %s", body)
	}

	var tokenResp dto.ExchangeCodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		log.Printf("failed to decode token response: %v", err)
		return model.GoogleUserInfo{}, fmt.Errorf("неверный ответ от сервера Google: %v", err)
	}

	userInfo, err := extractGoogleUserInfo(tokenResp.IdToken)
	if err != nil {
		return model.GoogleUserInfo{}, err
	}

	return userInfo, nil
}

func extractGoogleUserInfo(idTokenString string) (model.GoogleUserInfo, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(idTokenString, jwt.MapClaims{})
	if err != nil {
		log.Printf("failed to parse id token: %v", err)
		return model.GoogleUserInfo{}, fmt.Errorf("%w: %v", ErrParseIDToken, err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Printf("invalid claims type in id token")
		return model.GoogleUserInfo{}, fmt.Errorf("неверный формат ID Token claims")
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
		log.Printf("invalid id token: %s", errMsg)
		return model.GoogleUserInfo{}, fmt.Errorf("отсутствуют обязательные поля в ID Token: %s", errMsg)
	}

	userInfo := model.NewGoogleUserInfo(email, name, givenName, familyName, sub, picture, emailVerified)
	return userInfo, nil
}

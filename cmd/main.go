package main

import (
	"auth/config"
	api "auth/internal/api/http"
	_ "auth/internal/api/http/docs"
	"auth/internal/api/http/handlers"
	"auth/internal/application"
	"auth/internal/application/ports"
	"auth/internal/application/query"
	"auth/internal/infrastructure/oauth"
	"net/http"
	"time"
)

const ClientTimeout = 30 * time.Second

// @title Auth API
// @version 0.1.0
// @description API сервиса аутентификации
// @BasePath /
func main() {
	port := "6767"

	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	// Клиенты
	baseClient := http.Client{Timeout: ClientTimeout}
	oauthGoogleClient := oauth.NewGoogleClient(&cfg.GoogleAuth, &baseClient)
	oauthVkClient := oauth.NewVkClient(&cfg.VkAuth, &baseClient)

	// Фабрика клиентов
	clientFactory := oauth.NewClientFactory(map[application.OAuthProvider]ports.OAuthClient{
		application.GoogleProvider: oauthGoogleClient,
		application.VkProvider:     oauthVkClient,
	})

	//Инфраструктура
	oauthStateGenerator := oauth.NewOauthStateClient(cfg.OAuthSecretKey)

	// Команды и Запросы
	getOAuthRedirectQuery := query.NewGetOAuthRedirectUrlQuery(clientFactory)
	getOAuthCallbackQuery := query.NewGetOAuthCallbackQuery(clientFactory)

	// Handler
	oauthRedirectHandler := handlers.NewGetOAuthRedirectUrlHandler(getOAuthRedirectQuery, oauthStateGenerator)
	oauthGoogleCallbackHandler := handlers.NewGoogleOAuthCallbackHandler(getOAuthCallbackQuery, oauthStateGenerator)
	oauthVkCallbackHandler := handlers.NewVkOAuthCallbackHandler(getOAuthCallbackQuery, oauthStateGenerator)

	// Роутер
	router := api.NewRouter(oauthRedirectHandler, oauthGoogleCallbackHandler, oauthVkCallbackHandler)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		panic(err)
	}
}

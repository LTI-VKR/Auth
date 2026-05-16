package main

import (
	"auth/config"
	api "auth/internal/api/http"
	_ "auth/internal/api/http/docs"
	"auth/internal/api/http/handlers"
	"auth/internal/application/query"
	"auth/internal/infrastructure/oauth"
	"net/http"
)

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
	//
	//postgresPool, err := postgres.NewPool(cfg.DatabaseUrl)
	//if err != nil {
	//	panic("failed to create pool")
	//}
	//defer postgresPool.Close()
	//
	//MinioClient, err := minIO.NewMinioClient(cfg.Endpoint, cfg.Login, cfg.Password)
	//if err != nil {
	//	panic("failed to create minio client")
	//}

	// Клиенты
	oauthGoogleClient := oauth.NewOAuth2Client(&cfg.GoogleAuth)

	//Инфраструктура
	oauthStateGenerator := oauth.NewOauthStateClient(cfg.OAuthSecretKey)

	// Репозитории
	//cmdRepo := postgres.NewProfileCommandRepository(postgresPool)
	//qryRepo := postgres.NewProfileQueryRepository(postgresPool)
	//
	//minioRepo := minIO.NewAvatarMinioRepository(MinioClient, cfg.Bucket)
	//
	// Команды и Запросы
	getGoogleOAuthRedirectQuery := query.NewGetGoogleOAuthRedirectQuery(oauthGoogleClient)
	getGoogleOAuthCallbackQuery := query.NewGetGoogleOAuthCallbackQuery(oauthGoogleClient)
	//createCmd := command.NewCreateProfileCommand(cmdRepo)
	//getQry := query.NewGetProfileQuery(qryRepo)
	//getAvatarUploadUrlCmd := command.NewGetAvatarQuery(minioRepo)

	// Handler
	oauthGoogleRedirectHandler := handlers.NewGetGoogleOAuthRedirectHandler(getGoogleOAuthRedirectQuery, oauthStateGenerator)
	oauthGoogleCallbackHandler := handlers.NewGoogleOAuthCallbackHandler(getGoogleOAuthCallbackQuery, oauthStateGenerator)

	// Роутер
	router := api.NewRouter(oauthGoogleRedirectHandler, oauthGoogleCallbackHandler)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		panic(err)
	}
}

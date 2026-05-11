package main

import (
	"auth/config"
	api "auth/internal/api/http"
	_ "auth/internal/api/http/docs"
	"net/http"
)

// @title Auth API
// @version 0.1.0
// @description API сервиса аутентификации
// @BasePath /
func main() {
	port := "6767"

	_, err := config.NewConfig()
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

	// Репозитории
	//cmdRepo := postgres.NewProfileCommandRepository(postgresPool)
	//qryRepo := postgres.NewProfileQueryRepository(postgresPool)
	//
	//minioRepo := minIO.NewAvatarMinioRepository(MinioClient, cfg.Bucket)
	//
	//// Команды и Запросы
	//createCmd := command.NewCreateProfileCommand(cmdRepo)
	//updateCmd := command.NewUpdateProfileCommand(cmdRepo)
	//getQry := query.NewGetProfileQuery(qryRepo)
	//getAvatarUploadUrlCmd := command.NewGetAvatarQuery(minioRepo)
	//getAvatarDownloadUrlCmd := query.NewGetAvatarQuery(minioRepo)
	//
	//// Хендлеры
	//createHandler := handlers.NewCreateProfileHandler(createCmd)
	//getHandler := handlers.NewGetProfileHandler(getQry)
	//updateHandler := handlers.NewUpdateProfileHandler(updateCmd, getQry)
	//GetAvatarUploadUrlHandler := handlers.NewGetAvatarUploadUrlHandler(getAvatarUploadUrlCmd)
	//GetAvatarDownloadUrlHandler := handlers.NewGetAvatarDownloadUrlHandler(getAvatarDownloadUrlCmd)

	// Роутер
	router := api.NewRouter()

	if err := http.ListenAndServe(":"+port, router); err != nil {
		panic(err)
	}
}

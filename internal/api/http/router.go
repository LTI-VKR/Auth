package http

import (
	"auth/internal/api/http/handlers"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewRouter(
	getGoogleOAuthRedirectHandler *handlers.GetGoogleOAuthRedirectHandler,
	googleOAuthCallbackHandler *handlers.GoogleOAuthCallbackHandler,
) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(15 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:5174", "http://localhost:5175", "http://localhost:6767"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "Accept"},
		AllowCredentials: true,
	}))

	r.Get("/health", handlers.Health)
	r.Get("/auth/google/url", getGoogleOAuthRedirectHandler.Handle)
	r.Get("/auth/google/callback", googleOAuthCallbackHandler.Handle)

	r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json")))
	return r
}

package handlers

import (
	httperrors "auth/internal/api/http/errors"
	"auth/internal/application/ports"
	"auth/internal/application/query"
	"log"
	"net/http"
)

type GetGoogleOAuthRedirectHandler struct {
	query          *query.GetGoogleOAuthRedirectQuery
	stateGenerator ports.StateGenerator
}

func NewGetGoogleOAuthRedirectHandler(redirectQuery *query.GetGoogleOAuthRedirectQuery, stateGenerator ports.StateGenerator) *GetGoogleOAuthRedirectHandler {
	return &GetGoogleOAuthRedirectHandler{query: redirectQuery, stateGenerator: stateGenerator}
}

// Handle godoc
// @Summary Get google OAuth redirect
// @Tags OAuth
// @Produce json
// @Param return_to query string true "Ссылка на фронт, куда вернуть данные"
// @Success 302 {string} string "Redirect to Google OAuth"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/google/url [get]
func (h *GetGoogleOAuthRedirectHandler) Handle(w http.ResponseWriter, r *http.Request) {
	returnTo := r.URL.Query().Get("return_to")
	if returnTo == "" {
		httperrors.WriteError(w, r, httperrors.Map(httperrors.ErrMissingReturnTo))
		return
	}

	state, err := h.stateGenerator.GenerateState(returnTo)
	if err != nil {
		log.Printf("failed to generate state: %v", err)
		httperrors.WriteError(w, r, httperrors.Map(err))
		return
	}

	authUri, err := h.query.Execute(r.Context(), state)
	if err != nil {
		log.Printf("failed to get auth uri: %v", err)
		httperrors.WriteError(w, r, httperrors.Map(err))
		return
	}

	http.Redirect(w, r, authUri, http.StatusFound)
}

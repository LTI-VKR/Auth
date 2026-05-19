package handlers

import (
	httperrors "auth/internal/api/http/errors"
	"auth/internal/application"
	"auth/internal/application/ports"
	"auth/internal/application/query"
	"net/http"
)

type GetOAuthRedirectUrlHandler struct {
	query          *query.GetOAuthRedirectUrlQuery
	stateGenerator ports.StateGenerator
}

func NewGetOAuthRedirectUrlHandler(redirectQuery *query.GetOAuthRedirectUrlQuery, stateGenerator ports.StateGenerator) *GetOAuthRedirectUrlHandler {
	return &GetOAuthRedirectUrlHandler{query: redirectQuery, stateGenerator: stateGenerator}
}

// Handle godoc
// @Summary Get OAuth redirect
// @Tags OAuth
// @Produce json
// @Param return_to query string true "Ссылка на фронт, куда вернуть данные"
// @Param oauth_provider query string true "OAuth провайдер (google, vk)"
// @Success 302 {string} string "Redirect to OAuth provider"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/url [get]
func (h *GetOAuthRedirectUrlHandler) Handle(w http.ResponseWriter, r *http.Request) {
	returnTo := r.URL.Query().Get("return_to")
	if returnTo == "" {
		httperrors.WriteError(w, r, httperrors.Map(httperrors.ErrMissingReturnTo))
		return
	}

	providerStr := r.URL.Query().Get("oauth_provider")
	if providerStr == "" {
		httperrors.WriteError(w, r, httperrors.Map(httperrors.ErrInvalidOAuthProvider))
		return
	}

	provider := application.OAuthProvider(providerStr)
	if !provider.IsValid() {
		httperrors.WriteError(w, r, httperrors.Map(httperrors.ErrInvalidOAuthProvider))
		return
	}

	state, err := h.stateGenerator.GenerateState(returnTo)
	if err != nil {
		httperrors.WriteError(w, r, httperrors.Map(err))
		return
	}

	authUri, err := h.query.Execute(r.Context(), state, provider)
	if err != nil {
		httperrors.WriteError(w, r, httperrors.Map(err))
		return
	}

	http.Redirect(w, r, authUri, http.StatusFound)
}

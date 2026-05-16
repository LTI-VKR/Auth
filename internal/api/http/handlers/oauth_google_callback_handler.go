package handlers

import (
	"auth/internal/api/http/dto"
	httperrors "auth/internal/api/http/errors"
	"auth/internal/application/ports"
	"auth/internal/application/query"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type GoogleOAuthCallbackHandler struct {
	query          *query.GetGoogleOAuthCallbackQuery
	stateGenerator ports.StateGenerator
}

func NewGoogleOAuthCallbackHandler(query *query.GetGoogleOAuthCallbackQuery, stateGenerator ports.StateGenerator) *GoogleOAuthCallbackHandler {
	return &GoogleOAuthCallbackHandler{query: query, stateGenerator: stateGenerator}
}

// Handle godoc
// @Summary Google OAuth callback
// @Tags OAuth
// @Produce json
// @Param state query string true "OAuth state"
// @Param code query string true "Authorization code"
// @Param iss query string false "Issuer"
// @Param scope query string false "OAuth scope"
// @Param authuser query string false "Auth user"
// @Param prompt query string false "Prompt"
// @Success 302 {string} string "Redirect with user data"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/google/callback [get]
func (h *GoogleOAuthCallbackHandler) Handle(w http.ResponseWriter, r *http.Request) {
	urlQuery := r.URL.Query()
	params := dto.GoogleOAuthCallbackRequestDto{
		State:    urlQuery.Get("state"),
		Code:     urlQuery.Get("code"),
		Iss:      urlQuery.Get("iss"),
		Scope:    urlQuery.Get("scope"),
		Authuser: urlQuery.Get("authuser"),
		Prompt:   urlQuery.Get("prompt"),
	}

	if params.State == "" {
		log.Printf("missing state parameter in callback")
		httperrors.WriteError(w, r, httperrors.Map(httperrors.ErrMissingOAuthState))
		return
	}
	if params.Code == "" {
		log.Printf("missing code parameter in callback")
		httperrors.WriteError(w, r, httperrors.Map(httperrors.ErrMissingOAuthCode))
		return
	}

	returnTo, err := h.stateGenerator.VerifyAndExtractState(params.State)
	if err != nil {
		log.Printf("failed to verify state: %v", err)
		httperrors.WriteError(w, r, httperrors.Map(err))
		return
	}

	userInfo, err := h.query.Execute(r.Context(), params.Code)
	if err != nil {
		log.Printf("failed to exchange code: %v", err)
		httperrors.WriteError(w, r, httperrors.Map(err))
		return
	}

	redirectParams := url.Values{}
	redirectParams.Set("name", userInfo.Name)
	redirectParams.Set("email", userInfo.Email)
	redirectParams.Set("given_name", userInfo.GivenName)
	redirectParams.Set("family_name", userInfo.FamilyName)
	redirectParams.Set("sub", userInfo.SubId)
	redirectParams.Set("email_verified", strconv.FormatBool(userInfo.EmailVerified))
	redirectParams.Set("picture", userInfo.Picture)

	redirectURL := returnTo + "?" + redirectParams.Encode()
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

package handlers

import (
	"auth/internal/api/http/dto"
	httperrors "auth/internal/api/http/errors"
	"auth/internal/application"
	"auth/internal/application/model"
	"auth/internal/application/ports"
	"auth/internal/application/query"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type VkOAuthCallbackHandler struct {
	query          *query.GetOAuthCallbackQuery
	stateGenerator ports.StateGenerator
}

func NewVkOAuthCallbackHandler(query *query.GetOAuthCallbackQuery, stateGenerator ports.StateGenerator) *VkOAuthCallbackHandler {
	return &VkOAuthCallbackHandler{query: query, stateGenerator: stateGenerator}
}

// Handle godoc
// @Summary VK OAuth callback
// @Tags OAuth
// @Produce json
// @Param state query string true "OAuth state"
// @Param code query string true "Authorization code"
// @Param device_id query string true "VK device id"
// @Success 302 {string} string "Redirect with user data"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/vk/callback [get]
func (h *VkOAuthCallbackHandler) Handle(w http.ResponseWriter, r *http.Request) {
	urlQuery := r.URL.Query()
	params := dto.VkOAuthCallbackRequestDto{
		State:    urlQuery.Get("state"),
		Code:     urlQuery.Get("code"),
		DeviceId: urlQuery.Get("device_id"),
	}

	if params.State == "" {
		httperrors.WriteError(w, r, httperrors.Map(httperrors.ErrMissingOAuthState))
		return
	}
	if params.Code == "" {
		httperrors.WriteError(w, r, httperrors.Map(httperrors.ErrMissingOAuthCode))
		return
	}
	if params.DeviceId == "" {
		httperrors.WriteError(w, r, httperrors.Map(httperrors.ErrMissingOAuthDeviceId))
		return
	}

	returnTo, err := h.stateGenerator.VerifyAndExtractState(params.State)
	if err != nil {
		log.Printf("failed to verify state: %v", err)
		httperrors.WriteError(w, r, httperrors.Map(err))
		return
	}

	userInfo, err := h.query.Execute(r.Context(), model.ExchangeCodeParams{
		Code:     params.Code,
		State:    params.State,
		DeviceID: params.DeviceId,
	}, application.VkProvider)
	if err != nil {
		log.Printf("failed to exchange code: %v", err)
		httperrors.WriteError(w, r, httperrors.Map(err))
		return
	}

	redirectParams := url.Values{}
	redirectParams.Set("name", userInfo.Name)
	redirectParams.Set("email", userInfo.Email)
	redirectParams.Set("given_name", userInfo.FirstName)
	redirectParams.Set("family_name", userInfo.LastName)
	redirectParams.Set("sub", userInfo.Id)
	redirectParams.Set("email_verified", strconv.FormatBool(userInfo.EmailVerified))
	redirectParams.Set("picture", userInfo.AvatarURL)

	redirectURL := returnTo + "?" + redirectParams.Encode()
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

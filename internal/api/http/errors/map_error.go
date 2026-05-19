package errors

import (
	oauthinfra "auth/internal/infrastructure/oauth"
	stdErrors "errors"
	"net/http"
)

var (
	ErrInvalidJSON          = stdErrors.New("неверный формат JSON")
	ErrInvalidQueryParams   = stdErrors.New("неверные параметры запроса")
	ErrInvalidOAuthState    = stdErrors.New("состояние OAuth невалидно")
	ErrMissingOAuthState    = stdErrors.New("параметр state отсутствует")
	ErrMissingOAuthCode     = stdErrors.New("параметр code отсутствует")
	ErrMissingReturnTo      = stdErrors.New("параметр return_to отсутствует")
	ErrMissingOAuthDeviceId = stdErrors.New("параметр device_id отсутствует")
	ErrInvalidOAuthProvider = stdErrors.New("неподдерживаемый OAuth провайдер")
)

func Map(err error) MappedError {
	if err == nil {
		return NewBasicMapped(http.StatusInternalServerError, "INTERNAL", "внутренняя ошибка")
	}

	switch {
	case stdErrors.Is(err, ErrInvalidJSON):
		return NewBasicMapped(http.StatusBadRequest, "INVALID_JSON", "неверный формат JSON")
	case stdErrors.Is(err, ErrInvalidQueryParams):
		return NewBasicMapped(http.StatusBadRequest, "INVALID_QUERY_PARAMS", "неверные параметры запроса")
	case stdErrors.Is(err, ErrInvalidOAuthState), stdErrors.Is(err, oauthinfra.ErrInvalidState), stdErrors.Is(err, oauthinfra.ErrStateFormat):
		return NewBasicMapped(http.StatusBadRequest, "INVALID_STATE", "состояние OAuth невалидно")
	case stdErrors.Is(err, oauthinfra.ErrBadSignature):
		return NewBasicMapped(http.StatusBadRequest, "INVALID_STATE", "подпись состояния неверна")
	case stdErrors.Is(err, oauthinfra.ErrEmptyReturnTo), stdErrors.Is(err, ErrMissingReturnTo), stdErrors.Is(err, oauthinfra.ErrInvalidReturnTo):
		return NewBasicMapped(http.StatusBadRequest, "INVALID_RETURN_TO", "return_to невалиден")
	case stdErrors.Is(err, ErrMissingOAuthState):
		return NewBasicMapped(http.StatusBadRequest, "MISSING_PARAM", "параметр state обязателен")
	case stdErrors.Is(err, ErrMissingOAuthCode):
		return NewBasicMapped(http.StatusBadRequest, "MISSING_PARAM", "параметр code обязателен")
	case stdErrors.Is(err, ErrMissingOAuthDeviceId):
		return NewBasicMapped(http.StatusBadRequest, "MISSING_PARAM", "параметр device_id обязателен")
	case stdErrors.Is(err, oauthinfra.ErrUnmarshalState):
		return NewBasicMapped(http.StatusBadRequest, "INVALID_STATE", "не удалось разобрать состояние")
	case stdErrors.Is(err, oauthinfra.ErrOAuthCodeExchange), stdErrors.Is(err, oauthinfra.ErrTokenResponse), stdErrors.Is(err, oauthinfra.ErrParseIDToken), stdErrors.Is(err, oauthinfra.ErrInvalidIDToken), stdErrors.Is(err, oauthinfra.ErrMissingClaims):
		return NewBasicMapped(http.StatusInternalServerError, "INTERNAL", err.Error())
	case stdErrors.Is(err, ErrInvalidOAuthProvider), stdErrors.Is(err, oauthinfra.ErrInvalidOAuthProvider):
		return NewBasicMapped(http.StatusBadRequest, "INVALID_PROVIDER", "неподдерживаемый OAuth провайдер")
	default:
		return NewBasicMapped(http.StatusInternalServerError, "INTERNAL", err.Error())
	}
}

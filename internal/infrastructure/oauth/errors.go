package oauth

import stdErrors "errors"

// state
var (
	ErrInvalidState         = stdErrors.New("состояние OAuth невалидно")
	ErrBadSignature         = stdErrors.New("подпись состояния неверна")
	ErrStateFormat          = stdErrors.New("неправильный формат состояния")
	ErrGenerateNonce        = stdErrors.New("не удалось сгенерировать nonce")
	ErrUnmarshalState       = stdErrors.New("не удалось разобрать состояние")
	ErrEmptyReturnTo        = stdErrors.New("return_to не должен быть пустым")
	ErrInvalidReturnTo      = stdErrors.New("невалидный return_to URL")
	ErrInvalidOAuthProvider = stdErrors.New("неподдерживаемый OAuth провайдер")
)

// OAuth
var (
	ErrOAuthCodeExchange = stdErrors.New("ошибка при обмене кода на токены")
	ErrTokenResponse     = stdErrors.New("неверный ответ от сервера Google")
	ErrParseIDToken      = stdErrors.New("не удалось разобрать ID Token")
	ErrInvalidIDToken    = stdErrors.New("неверный формат ID Token claims")
	ErrMissingClaims     = stdErrors.New("отсутствуют обязательные поля в ID Token")
)

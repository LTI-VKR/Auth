package errors

import stdErrors "errors"

var (
	ErrInvalidJSON        = stdErrors.New("неверный формат JSON")
	ErrInvalidQueryParams = stdErrors.New("неверные параметры запроса")
	ErrInvalidOAuthState  = stdErrors.New("состояние OAuth невалидно")
	ErrMissingOAuthState  = stdErrors.New("параметр state отсутствует")
	ErrMissingOAuthCode   = stdErrors.New("параметр code отсутствует")
	ErrMissingReturnTo    = stdErrors.New("параметр return_to отсутствует")
)

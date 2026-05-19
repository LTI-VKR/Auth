package oauth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/url"
	"strings"
)

type StateGenerator struct {
	secretKey []byte
}

type stateData struct {
	ReturnTo string `json:"return_to"`
	Nonce    string `json:"nonce"`
}

func NewOauthStateClient(SecretKey string) *StateGenerator {
	return &StateGenerator{
		secretKey: []byte(SecretKey),
	}
}

func (s *StateGenerator) GenerateState(returnTo string) (string, error) {
	// Валидация returnTo
	if returnTo == "" {
		return "", ErrEmptyReturnTo
	}

	// Проверка, что это валидный URL
	if _, err := url.Parse(returnTo); err != nil {
		return "", ErrInvalidReturnTo
	}

	nonce := make([]byte, 16)
	_, err := rand.Read(nonce)
	if err != nil {
		return "", ErrGenerateNonce
	}

	data := &stateData{
		ReturnTo: returnTo,
		Nonce:    base64.RawURLEncoding.EncodeToString(nonce),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	mac := hmac.New(sha256.New, s.secretKey)
	mac.Write(jsonData)
	sig := mac.Sum(nil)

	encoded := base64.RawURLEncoding.EncodeToString(jsonData)
	sigEncoded := base64.RawURLEncoding.EncodeToString(sig)

	state := encoded + "." + sigEncoded
	// Encode full state to avoid '.' being stripped by providers.
	return base64.RawURLEncoding.EncodeToString([]byte(state)), nil
}

func (s *StateGenerator) VerifyAndExtractState(state string) (string, error) {
	decodedState, err := base64.RawURLEncoding.DecodeString(state)
	if err != nil {
		return "", ErrInvalidState
	}

	parts := strings.Split(string(decodedState), ".")
	if len(parts) != 2 {
		return "", ErrStateFormat
	}

	jsonData, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return "", ErrInvalidState
	}

	sig, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", ErrInvalidState
	}

	mac := hmac.New(sha256.New, s.secretKey)
	mac.Write(jsonData)
	expected := mac.Sum(nil)

	if !hmac.Equal(sig, expected) {
		return "", ErrBadSignature
	}

	var data stateData
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return "", ErrUnmarshalState
	}

	if data.ReturnTo == "" {
		return "", ErrEmptyReturnTo
	}

	return data.ReturnTo, nil
}

package model

type ExchangeCodeParams struct {
	Code         string
	State        string
	DeviceID     string
	CodeVerifier string
}

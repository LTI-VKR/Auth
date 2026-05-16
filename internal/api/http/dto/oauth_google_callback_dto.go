package dto

type GoogleOAuthCallbackRequestDto struct {
	State    string `json:"state"`
	Code     string `json:"code"`
	Iss      string `json:"iss,omitempty"`
	Scope    string `json:"scope,omitempty"`
	Authuser string `json:"authuser,omitempty"`
	Prompt   string `json:"prompt,omitempty"`
}

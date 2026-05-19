package dto

type VkOAuthCallbackRequestDto struct {
	State    string `json:"state"`
	Code     string `json:"code"`
	DeviceId string `json:"device_id"`
}

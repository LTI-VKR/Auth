package dto

type ErrorResponse struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
}

func (e *ErrorResponse) SetRequestID(id string) {
	e.RequestID = id
}

type ErrorValidationResponse struct {
	Code      string            `json:"code"`
	Fields    map[string]string `json:"fields"`
	RequestID string            `json:"request_id"`
}

func (e *ErrorValidationResponse) SetRequestID(id string) {
	e.RequestID = id
}

package auth

type ValidateTOPTCodeResponse struct {
	Reason  string `json:"reason,omitempty"`
	IsValid bool   `json:"isValid"`
}

type GenerateTOPTURIResponse struct {
	URI string `json:"uri"`
}

type LoginResponse struct {
	Success bool    `json:"success"`
	Reason  string  `json:"reason,omitempty"`
	Token   *string `json:"token"`
}

type AuthMiddlewareResponse struct {
	Success bool   `json:"success"`
	Reason  string `json:"reason,omitempty"`
}

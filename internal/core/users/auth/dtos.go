package auth

type ValidateTOPTCodeResponse struct {
	IsValid    bool   `json:"isValid"`
	Reason     string `json:"reason,omitempty"`
	ReasonCode string `json:"reasonCode,omitempty"`
}

type GenerateTOPTURIResponse struct {
	URI string `json:"uri"`
}

type LoginResponse struct {
	Success    bool    `json:"success"`
	Token      *string `json:"token"`
	Reason     string  `json:"reason,omitempty"`
	ReasonCode string  `json:"reasonCode,omitempty"`
}

type AuthMiddlewareResponse struct {
	Success    bool   `json:"success"`
	Reason     string `json:"reason,omitempty"`
	ReasonCode string `json:"reasonCode,omitempty"`
}

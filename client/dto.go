package client

type PaymentRequest struct {
	AccountId string  `json:"account_id"`
	Reference string  `json:"reference"`
	Amount    float64 `json:"amount"`
}

type ErrorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

type PaymentResponse struct {
	AccountId string  `json:"account_id"`
	Reference string  `json:"reference"`
	Amount    float64 `json:"amount"`
}

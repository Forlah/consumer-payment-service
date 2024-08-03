package models

type PaymentRequestPayload struct {
	UserId    string  `json:"user_id"`
	AccountId string  `json:"account_id"`
	Reference string  `json:"reference"`
	Amount    float64 `json:"amount"`
}

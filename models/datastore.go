package models

type User struct {
	Id        string `bson:"id"`
	AccountID string `bson:"account_id"`
	Name      string `bson:"account"`
	CreatedAt int64  `bson:"created_at"`
}

type Account struct {
	AccountID string  `bson:"account_id"`
	Balance   float64 `bson:"balance"`
	CreatedAt int64   `bson:"created_at"`
}

type TransactionType string

const (
	DEBIT  TransactionType = "DEBIT"
	CREDIT TransactionType = "CREDIT"
)

type TransactionStatus string

const (
	SUCCESS TransactionStatus = "SUCCESS"
	FAILED  TransactionStatus = "FAILED"
)

type Transaction struct {
	Reference string            `bson:"reference"`
	UserID    string            `bson:"user_id"`
	AccountID string            `bson:"account_id"`
	Amount    float64           `bson:"amount"`
	Type      TransactionType   `bson:"type"`
	Status    TransactionStatus `bson:"status"`
	CreatedAt int64             `bson:"created_at"`
}

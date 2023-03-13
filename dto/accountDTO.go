package dto

type TransferRequest struct {
	DebitAccountId  string  `json:"debitAccountId" validate:"required"`
	CreditAccountId string  `json:"creditAccountId" validate:"required"`
	Amount          float64 `json:"amount" validate:"required,gte=0"`
}

type TransferResponse struct {
	Success bool `json:"success"`
}

type AccountsResponse struct {
	Accounts []Account `json:"accounts"`
}

type Account struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
}

type Error struct {
	Item    string `json:"items"`
	Message string `json:"message"`
}

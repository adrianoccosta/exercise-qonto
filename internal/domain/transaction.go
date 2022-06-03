package domain

// Transaction Struct that represents a payment transaction
type Transaction struct {
	Name             string  `json:"name"`
	Iban             string  `json:"iban"`
	Bic              string  `json:"bic"`
	CounterPartyName string  `json:"counterparty_name"`
	CounterPartyIban string  `json:"counterparty_iban"`
	CounterPartyBic  string  `json:"counterparty_bic"`
	Amount           float64 `json:"amount,string"`
	Currency         string  `json:"currency"`
	Description      string  `json:"description"`
}

// TransactionList Struct that represents a list of payment transactions
type TransactionList []Transaction

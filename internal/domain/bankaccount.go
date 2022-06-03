package domain

// BankAccount Struct that represents a use back account
type BankAccount struct {
	Name    string  `json:"name"`
	Balance float64 `json:"balance,string"`
	Iban    string  `json:"iban"`
	Bic     string  `json:"bic"`
}

package domain

import "github.com/go-playground/validator/v10"

// BankAccount Struct that represents a use back account
type BankAccount struct {
	Name    string  `json:"name" validate:"required"`
	Balance float64 `json:"balance,string" validate:"required"`
	Iban    string  `json:"iban" validate:"required"`
	Bic     string  `json:"bic" validate:"required"`
}

//Validate validates the BankAccount struct based on 'validate' tags of its fields
func (l *BankAccount) Validate() error {
	v := validator.New()
	return v.Struct(l)
}

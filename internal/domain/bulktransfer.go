package domain

import "github.com/go-playground/validator/v10"

// BulkTransfer Struct that represents a payment request
type BulkTransfer struct {
	OrganizationName string           `json:"organization_name" validate:"required"`
	OrganizationBic  string           `json:"organization_bic" validate:"required"`
	OrganizationIban string           `json:"organization_iban" validate:"required"`
	CreditTransfers  []CreditTransfer `json:"credit_transfers" validate:"required"`
}

type CreditTransfer struct {
	Amount           float64 `json:"amount,string" validate:"required"`
	Currency         string  `json:"currency" validate:"required"`
	CounterPartyName string  `json:"counterparty_name" validate:"required"`
	CounterPartyBic  string  `json:"counterparty_bic" validate:"required"`
	CounterPartyIban string  `json:"counterparty_iban" validate:"required"`
	Description      string  `json:"description"`
}

//Validate validates the BulkTransfer struct based on 'validate' tags of its fields
func (l *BulkTransfer) Validate() error {
	v := validator.New()
	return v.Struct(l)
}

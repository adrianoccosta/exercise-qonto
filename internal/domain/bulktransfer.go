package domain

// BulkTransfer Struct that represents a payment request
type BulkTransfer struct {
	OrganizationName string           `json:"organization_name"`
	OrganizationBic  string           `json:"organization_bic"`
	OrganizationIban string           `json:"organization_iban"`
	CreditTransfers  []CreditTransfer `json:"credit_transfers"`
}

type CreditTransfer struct {
	Amount           float64 `json:"amount,string"`
	Currency         string  `json:"currency"`
	CounterPartyName string  `json:"counterparty_name"`
	CounterPartyBic  string  `json:"counterparty_bic"`
	CounterPartyIban string  `json:"counterparty_iban"`
	Description      string  `json:"description"`
}

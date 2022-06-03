package bankaccountrepo

import (
	"errors"
	"github.com/adrianoccosta/exercise-qonto/cmd/config"
)

// Repo struct
type Repo struct {
	DB config.Conn
}

// BankAccount Struct that represents a use back account
type BankAccount struct {
	ID               uint
	OrganizationName string
	BalanceCents     int
	Iban             string
	Bic              string
}

// BankAccountRepository Interface for the back account registry
type BankAccountRepository interface {
	Create(data BankAccount) (int, error)
	Read(bankAccountID uint) (BankAccount, error)
	ReadByIban(iban string) (BankAccount, error)
	Update(data BankAccount) error
	DeleteByIban(iban string) error
}

// New Returns a new instance of DB.
func New(db config.Conn) Repo {
	return Repo{
		DB: db,
	}
}

// Create new bank account
func (repo Repo) Create(data BankAccount) (int, error) {
	if res, _ := repo.ReadByIban(data.Iban); (res != BankAccount{}) {
		return 0, errors.New("Register with same iban already exists")
	}
	insertQuery := "INSERT INTO bank_accounts" +
		"(organization_name, balance_cents, iban, bic) " +
		"VALUES (?, ?, ?, ?)"

	res, err := repo.DB.Conn.Exec(insertQuery, data.OrganizationName, data.BalanceCents, data.Iban, data.Bic)

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()

	return int(id), err
}

// Read a bank account
func (repo Repo) Read(bankAccountID uint) (BankAccount, error) {
	query := "SELECT id, organization_name, balance_cents, iban, bic " +
		" FROM bank_accounts" +
		" WHERE id = ?"

	row := repo.DB.Conn.QueryRow(query, bankAccountID)

	var bankAccount BankAccount
	err := row.Scan(
		&bankAccount.ID,
		&bankAccount.OrganizationName,
		&bankAccount.BalanceCents,
		&bankAccount.Iban,
		&bankAccount.Bic,
	)
	if err != nil {
		return BankAccount{}, err
	}

	return bankAccount, nil
}

// ReadByIban a bank account
func (repo Repo) ReadByIban(iban string) (BankAccount, error) {
	query := "SELECT id, organization_name, balance_cents, iban, bic " +
		" FROM bank_accounts" +
		" WHERE iban = ?"

	row := repo.DB.Conn.QueryRow(query, iban)

	var bankAccount BankAccount
	err := row.Scan(
		&bankAccount.ID,
		&bankAccount.OrganizationName,
		&bankAccount.BalanceCents,
		&bankAccount.Iban,
		&bankAccount.Bic,
	)
	if err != nil {
		return BankAccount{}, err
	}

	return bankAccount, nil
}

// Update the values of a bank account
func (repo Repo) Update(data BankAccount) error {
	updateQuery := "UPDATE bank_accounts " +
		"SET organization_name = ?, balance_cents = ?, bic = ? " +
		"WHERE id = ?"

	_, err := repo.DB.Conn.Exec(updateQuery, data.OrganizationName, data.BalanceCents, data.Bic, data.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteByIban a bank account
func (repo Repo) DeleteByIban(iban string) error {
	deleteQuery := "DELETE " +
		" FROM bank_accounts" +
		" WHERE iban = ?"

	_, err := repo.DB.Conn.Exec(deleteQuery, iban)

	if err != nil {
		return err
	}

	return nil
}

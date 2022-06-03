package transactionrepo

import (
	"database/sql"
	"fmt"
	"github.com/adrianoccosta/exercise-qonto/cmd/config"
	"github.com/adrianoccosta/exercise-qonto/internal/domain"
)

// Repo struct
type Repo struct {
	DB config.Conn
}

// TransactionList list of Transaction
type TransactionList []Transaction

// Transaction Struct that represents a payment transaction
type Transaction struct {
	ID               uint
	CounterPartyName string
	CounterPartyIban string
	CounterPartyBic  string
	AmountCents      int
	AmountCurrency   string
	BankAccountID    uint
	Description      string
}

// TransactionRepository Interface for the payment transactions
type TransactionRepository interface {
	Create(data Transaction) (int, error)
	Read(transactionID uint) (Transaction, error)
	ReadByFilter(filters map[string]string) (domain.TransactionList, error)
}

// New Returns a new instance of DB.
func New(db config.Conn) Repo {
	return Repo{
		DB: db,
	}
}

// Create new transaction
func (repo Repo) Create(data Transaction) (int, error) {
	insertQuery := "INSERT INTO transactions" +
		"(counterparty_name, counterparty_iban, counterparty_bic, amount_cents, amount_currency, bank_account_id, description) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?)"

	res, err := repo.DB.Conn.Exec(
		insertQuery,
		data.CounterPartyName,
		data.CounterPartyIban,
		data.CounterPartyBic,
		data.AmountCents,
		data.AmountCurrency,
		data.BankAccountID,
		data.Description)

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()

	return int(id), err
}

// Read a transaction
func (repo Repo) Read(transactionID uint) (Transaction, error) {
	query := "SELECT id, counterparty_name, counterparty_iban, counterparty_bic, amount_cents, amount_currency, bank_account_id, description " +
		" FROM transactions" +
		" WHERE id = ?"

	row := repo.DB.Conn.QueryRow(query, transactionID)

	var transaction Transaction
	err := row.Scan(
		&transaction.ID,
		&transaction.CounterPartyName,
		&transaction.CounterPartyIban,
		&transaction.CounterPartyBic,
		&transaction.AmountCents,
		&transaction.AmountCurrency,
		&transaction.BankAccountID,
		&transaction.Description,
	)
	if err != nil {
		return Transaction{}, err
	}

	return transaction, nil
}

// ReadByFilter a transaction
func (repo Repo) ReadByFilter(filters map[string]string) (domain.TransactionList, error) {
	query := "SELECT organization_name as name, iban, bic, counterparty_name, counterparty_iban, counterparty_bic, CAST(amount_cents AS float)/100 as amount, amount_currency, description" +
		" FROM transactions" +
		" INNER JOIN bank_accounts b ON b.id = bank_account_id" +
		" WHERE 1 = 1"

	var bind []any
	for k, v := range filters {
		query += fmt.Sprintf(" and %s = ?", k)
		bind = append(bind, v)
	}

	rows, err := repo.DB.Conn.Query(query, bind...)

	if err != nil {
		return nil, err
	}

	return createFromDB(rows)
}

// createFromDB Transaction List mapper
func createFromDB(rows *sql.Rows) (domain.TransactionList, error) {
	var allTransactions domain.TransactionList
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	for rows.Next() {
		var transaction domain.Transaction
		err := rows.Scan(
			&transaction.Name,
			&transaction.Iban,
			&transaction.Bic,
			&transaction.CounterPartyName,
			&transaction.CounterPartyIban,
			&transaction.CounterPartyBic,
			&transaction.Amount,
			&transaction.Currency,
			&transaction.Description,
		)

		if err != nil {
			return nil, err
		}

		allTransactions = append(allTransactions, transaction)
	}

	return allTransactions, nil
}

package transactionrepo

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/adrianoccosta/exercise-qonto/cmd/config"
	"github.com/adrianoccosta/exercise-qonto/internal/repository/bankaccountrepo"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func setupTransactionRepo() (config.Conn, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return config.Conn{Conn: db}, mock
}

func TestTransactionRepo(t *testing.T) {

	conn, mock := setupTransactionRepo()
	defer func() {
		mock.ExpectClose()
		err := conn.Conn.Close()
		if err != nil {
			t.Errorf("Error closing connection: %+v", err)
		}
	}()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := Repo{DB: conn}

	t.Run("Test constructor.", func(t *testing.T) {
		r := New(conn)

		assert.NotEmpty(t, r)
	})

	transaction := Transaction{
		ID:               22,
		CounterPartyName: "Bip Bip",
		CounterPartyIban: "EE383680981021245685",
		CounterPartyBic:  "CRLYFRPPTOU",
		AmountCents:      123123,
		AmountCurrency:   "EUR",
		BankAccountID:    1,
		Description:      "Wonderland/4410",
	}

	bankAccount := bankaccountrepo.BankAccount{
		ID:               1,
		OrganizationName: "ACME Corp",
		BalanceCents:     123456,
		Iban:             "FR10474608000002006107XXXXX",
		Bic:              "OIVUSCLQXXX",
	}

	t.Run("Test Create return success", func(t *testing.T) {
		insertQuery := "INSERT INTO transactions"

		mock.ExpectExec(insertQuery).
			WithArgs(transaction.CounterPartyName,
				transaction.CounterPartyIban,
				transaction.CounterPartyBic,
				transaction.AmountCents,
				transaction.AmountCurrency,
				transaction.BankAccountID,
				transaction.Description).
			WillReturnResult(sqlmock.NewResult(1, 1))

		r, err := repo.Create(transaction)
		assert.NoError(t, err)
		assert.Equal(t, 1, r)
	})

	t.Run("Test Create return error while inserting on database.", func(t *testing.T) {
		insertQuery := "INSERT INTO transactions"

		mock.ExpectExec(insertQuery).
			WithArgs(transaction.CounterPartyName,
				transaction.CounterPartyIban,
				transaction.CounterPartyBic,
				transaction.AmountCents,
				transaction.AmountCurrency,
				transaction.BankAccountID,
				transaction.Description).
			WillReturnError(fmt.Errorf("error"))

		_, err := repo.Create(transaction)
		assert.Error(t, err)
	})

	t.Run("Test Read return success", func(t *testing.T) {
		selectQuery := "SELECT id, counterparty_name, counterparty_iban, counterparty_bic, amount_cents, amount_currency, bank_account_id, description"

		rows := sqlmock.NewRows([]string{"id", "counterparty_name", "counterparty_iban", "counterparty_bic", "amount_cents", "amount_currency", "bank_account_id", "description"})
		rows.AddRow(transaction.ID, transaction.CounterPartyName, transaction.CounterPartyIban, transaction.CounterPartyBic, transaction.AmountCents, transaction.AmountCurrency, transaction.BankAccountID, transaction.Description)

		mock.ExpectQuery(selectQuery).
			WithArgs(transaction.ID).
			WillReturnRows(rows)

		s, err := repo.Read(transaction.ID)
		assert.NoError(t, err)

		assert.Equal(t, transaction.ID, s.ID)
		assert.Equal(t, transaction.CounterPartyName, s.CounterPartyName)
		assert.Equal(t, transaction.CounterPartyIban, s.CounterPartyIban)
		assert.Equal(t, transaction.CounterPartyBic, s.CounterPartyBic)
		assert.Equal(t, transaction.AmountCents, s.AmountCents)
		assert.Equal(t, transaction.AmountCurrency, s.AmountCurrency)
		assert.Equal(t, transaction.BankAccountID, s.BankAccountID)
		assert.Equal(t, transaction.Description, s.Description)
	})

	t.Run("Test Read return error", func(t *testing.T) {
		selectQuery := "SELECT id, counterparty_name, counterparty_iban, counterparty_bic, amount_cents, amount_currency, bank_account_id, description"

		mock.ExpectQuery(selectQuery).
			WithArgs(transaction.ID).
			WillReturnError(fmt.Errorf("error"))

		_, err := repo.Read(transaction.ID)
		assert.Error(t, err)
	})

	t.Run("Test ReadByFilter return success", func(t *testing.T) {
		selectQuery := "SELECT organization_name as name, iban, bic, counterparty_name, counterparty_iban, counterparty_bic,"

		rows := sqlmock.NewRows([]string{"organization_name", "iban", "bic", "counterparty_name", "counterparty_iban", "counterparty_bic", "amount_cents", "amount_currency", "description"})
		rows.AddRow(bankAccount.OrganizationName, bankAccount.Iban, bankAccount.Bic, transaction.CounterPartyName, transaction.CounterPartyIban, transaction.CounterPartyBic, float64(transaction.AmountCents)/100, transaction.AmountCurrency, transaction.Description)

		mock.ExpectQuery(selectQuery).
			WillReturnRows(rows)

		var filters map[string]string

		s, err := repo.ReadByFilter(filters)
		assert.NoError(t, err)
		assert.Equal(t, bankAccount.OrganizationName, s[0].Name)
		assert.Equal(t, bankAccount.Iban, s[0].Iban)
		assert.Equal(t, bankAccount.Bic, s[0].Bic)
		assert.Equal(t, transaction.CounterPartyName, s[0].CounterPartyName)
		assert.Equal(t, transaction.CounterPartyIban, s[0].CounterPartyIban)
		assert.Equal(t, transaction.CounterPartyBic, s[0].CounterPartyBic)
		assert.Equal(t, 1231.23, s[0].Amount)
		assert.Equal(t, transaction.AmountCurrency, s[0].Currency)
		assert.Equal(t, transaction.Description, s[0].Description)
	})
}

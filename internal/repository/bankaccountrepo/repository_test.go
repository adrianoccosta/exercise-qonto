package bankaccountrepo

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/adrianoccosta/exercise-qonto/cmd/config"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func setupBankAccountRepo() (config.Conn, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return config.Conn{Conn: db}, mock
}

func TestBankAccountRepo(t *testing.T) {

	conn, mock := setupBankAccountRepo()
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

	bankAccount := BankAccount{
		ID:               1,
		OrganizationName: "ACME Corp",
		BalanceCents:     123456,
		Iban:             "FR10474608000002006107XXXXX",
		Bic:              "OIVUSCLQXXX",
	}

	t.Run("Test Create return success", func(t *testing.T) {
		insertQuery := "INSERT INTO bank_accounts"

		mock.ExpectExec(insertQuery).
			WithArgs(bankAccount.OrganizationName,
				bankAccount.BalanceCents,
				bankAccount.Iban,
				bankAccount.Bic).
			WillReturnResult(sqlmock.NewResult(1, 1))

		r, err := repo.Create(bankAccount)
		assert.NoError(t, err)
		assert.Equal(t, 1, r)
	})

	t.Run("Test Create return error while inserting on database.", func(t *testing.T) {
		insertQuery := "INSERT INTO bank_accounts"

		mock.ExpectExec(insertQuery).
			WithArgs(bankAccount.OrganizationName,
				bankAccount.BalanceCents,
				bankAccount.Iban,
				bankAccount.Bic).
			WillReturnError(fmt.Errorf("error"))

		_, err := repo.Create(bankAccount)
		assert.Error(t, err)
	})

	t.Run("Test Read return success", func(t *testing.T) {
		selectQuery := "SELECT id, organization_name, balance_cents, iban, bic FROM bank_accounts"

		rows := sqlmock.NewRows([]string{"id", "organization_name", "balance_cents", "iban", "bic"})
		rows.AddRow(bankAccount.ID, bankAccount.OrganizationName, bankAccount.BalanceCents, bankAccount.Iban, bankAccount.Bic)

		mock.ExpectQuery(selectQuery).
			WithArgs(bankAccount.ID).
			WillReturnRows(rows)

		s, err := repo.Read(bankAccount.ID)
		assert.NoError(t, err)

		assert.Equal(t, bankAccount.ID, s.ID)
		assert.Equal(t, bankAccount.OrganizationName, s.OrganizationName)
		assert.Equal(t, bankAccount.BalanceCents, s.BalanceCents)
		assert.Equal(t, bankAccount.Iban, s.Iban)
		assert.Equal(t, bankAccount.Bic, s.Bic)
	})

	t.Run("Test Read return error", func(t *testing.T) {
		selectQuery := "SELECT id, organization_name, balance_cents, iban, bic FROM bank_accounts"

		mock.ExpectQuery(selectQuery).
			WithArgs(bankAccount.ID).
			WillReturnError(fmt.Errorf("error"))

		_, err := repo.Read(bankAccount.ID)
		assert.Error(t, err)
	})

	t.Run("Test ReadByIban return success", func(t *testing.T) {
		selectQuery := "SELECT id, organization_name, balance_cents, iban, bic FROM bank_accounts"

		rows := sqlmock.NewRows([]string{"id", "organization_name", "balance_cents", "iban", "bic"})
		rows.AddRow(bankAccount.ID, bankAccount.OrganizationName, bankAccount.BalanceCents, bankAccount.Iban, bankAccount.Bic)

		mock.ExpectQuery(selectQuery).
			WithArgs(bankAccount.Iban).
			WillReturnRows(rows)

		s, err := repo.ReadByIban(bankAccount.Iban)
		assert.NoError(t, err)

		assert.Equal(t, bankAccount.ID, s.ID)
		assert.Equal(t, bankAccount.OrganizationName, s.OrganizationName)
		assert.Equal(t, bankAccount.BalanceCents, s.BalanceCents)
		assert.Equal(t, bankAccount.Iban, s.Iban)
		assert.Equal(t, bankAccount.Bic, s.Bic)
	})

	t.Run("Test ReadByIban return error", func(t *testing.T) {
		selectQuery := "SELECT id, organization_name, balance_cents, iban, bic FROM bank_accounts"

		mock.ExpectQuery(selectQuery).
			WithArgs(bankAccount.Iban).
			WillReturnError(fmt.Errorf("error"))

		_, err := repo.ReadByIban(bankAccount.Iban)
		assert.Error(t, err)
	})

	t.Run("Test Update return success.", func(t *testing.T) {

		updateQuery := "UPDATE bank_accounts"

		mock.ExpectExec(updateQuery).
			WithArgs(bankAccount.OrganizationName, bankAccount.BalanceCents, bankAccount.Bic, bankAccount.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Update(bankAccount)
		assert.NoError(t, err)
	})

	t.Run("Test Update return error.", func(t *testing.T) {

		updateQuery := "UPDATE bank_accounts"

		mock.ExpectExec(updateQuery).
			WithArgs(bankAccount.OrganizationName, bankAccount.BalanceCents, bankAccount.Bic, bankAccount.ID).
			WillReturnError(fmt.Errorf("error"))

		err := repo.Update(bankAccount)
		assert.Error(t, err)
	})

	t.Run("Test DeleteByIban return success.", func(t *testing.T) {

		deleteQuery := "DELETE"

		mock.ExpectExec(deleteQuery).
			WithArgs(bankAccount.Iban).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.DeleteByIban(bankAccount.Iban)
		assert.NoError(t, err)
	})
}

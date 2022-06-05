package bankaccountsvc

import (
	"errors"
	"github.com/adrianoccosta/exercise-qonto/internal/domain"
	"github.com/adrianoccosta/exercise-qonto/internal/repository/bankaccountrepo"
	"github.com/adrianoccosta/exercise-qonto/test/mocks/log"
	"github.com/adrianoccosta/exercise-qonto/test/mocks/repository"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBankAccountService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := mockrepository.NewMockBankAccountRepository(ctrl)
	logMock := mocklog.NewMockLogger(ctrl)

	bankAccount := domain.BankAccount{
		Name:    "ACME Corp",
		Balance: 12.40,
		Iban:    "FR10474608000002006107XXXXX",
		Bic:     "OIVUSCLQXXX",
	}

	bankAccountRepo := bankaccountrepo.BankAccount{
		OrganizationName: "ACME Corp",
		BalanceCents:     1240,
		Iban:             "FR10474608000002006107XXXXX",
		Bic:              "OIVUSCLQXXX",
	}

	t.Run("Test Create return success", func(t *testing.T) {
		repoMock.EXPECT().
			Create(bankAccountRepo).
			Return(1, nil)

		svc := New(repoMock, logMock)
		err := svc.Create(bankAccount)

		assert.Nil(t, err)
	})

	t.Run("Test Create return error", func(t *testing.T) {
		repoMock.EXPECT().
			Create(gomock.Any()).
			Return(0, errors.New("error"))

		svc := New(repoMock, logMock)
		err := svc.Create(domain.BankAccount{})

		assert.Error(t, err)
	})

	t.Run("Test Read return success", func(t *testing.T) {
		repoMock.EXPECT().
			ReadByIban("FR10474608000002006107XXXXX").
			Return(bankAccountRepo, nil)

		svc := New(repoMock, logMock)
		res, err := svc.Read("FR10474608000002006107XXXXX")

		assert.Nil(t, err)
		assert.Equal(t, bankAccount, res)
	})

	t.Run("Test Read return error", func(t *testing.T) {
		repoMock.EXPECT().
			ReadByIban(gomock.Any()).
			Return(bankaccountrepo.BankAccount{}, errors.New("error"))

		svc := New(repoMock, logMock)
		_, err := svc.Read("FR10474608000002006107XXXXX")

		assert.Error(t, err)
	})

	t.Run("Test Update return success", func(t *testing.T) {
		repoMock.EXPECT().
			ReadByIban("FR10474608000002006107XXXXX").
			Return(bankAccountRepo, nil)
		repoMock.EXPECT().
			Update(gomock.Any()).
			Return(nil)

		svc := New(repoMock, logMock)
		err := svc.Update(bankAccount)

		assert.Nil(t, err)
	})

	t.Run("Test Update return error reading value", func(t *testing.T) {
		repoMock.EXPECT().
			ReadByIban(gomock.Any()).
			Return(bankaccountrepo.BankAccount{}, errors.New("error"))

		svc := New(repoMock, logMock)
		err := svc.Update(domain.BankAccount{})

		assert.Error(t, err)
	})

	t.Run("Test Update return error updating value", func(t *testing.T) {
		repoMock.EXPECT().
			ReadByIban(gomock.Any()).
			Return(bankAccountRepo, nil)
		repoMock.EXPECT().
			Update(gomock.Any()).
			Return(errors.New("error"))

		svc := New(repoMock, logMock)
		err := svc.Update(domain.BankAccount{})

		assert.Error(t, err)
	})

	t.Run("Test Delete return success", func(t *testing.T) {
		repoMock.EXPECT().
			DeleteByIban("FR10474608000002006107XXXXX").
			Return(nil)

		svc := New(repoMock, logMock)
		err := svc.Delete("FR10474608000002006107XXXXX")

		assert.Nil(t, err)
	})

	t.Run("Test Delete return error", func(t *testing.T) {
		repoMock.EXPECT().
			DeleteByIban(gomock.Any()).
			Return(errors.New("error"))

		svc := New(repoMock, logMock)
		err := svc.Delete("FR10474608000002006107XXXXX")

		assert.Error(t, err)
	})
}

package transfersvc

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

func TestTransferService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMockTransaction := mockrepository.NewMockTransactionRepository(ctrl)
	repoMockBankAccount := mockrepository.NewMockBankAccountRepository(ctrl)
	logMock := mocklog.NewMockLogger(ctrl)

	bankAccountRepo := bankaccountrepo.BankAccount{
		ID:               1,
		OrganizationName: "ACME Corp",
		BalanceCents:     1460,
		Iban:             "FR10474608000002006107XXXXX",
		Bic:              "OIVUSCLQXXX",
	}

	bulkTransfer := domain.BulkTransfer{
		OrganizationName: "ACME Corp",
		OrganizationBic:  "OIVUSCLQXXX",
		OrganizationIban: "FR10474608000002006107XXXXX",
		CreditTransfers: []domain.CreditTransfer{
			{
				Amount:           14.53,
				Currency:         "EUR",
				CounterPartyName: "Bip Bip",
				CounterPartyBic:  "CRLYFRPPTOU",
				CounterPartyIban: "EE383680981021245685",
				Description:      "Wonderland/4410",
			},
		},
	}

	t.Run("Test BulkTransfer return success", func(t *testing.T) {
		repoMockBankAccount.EXPECT().
			ReadByIban("FR10474608000002006107XXXXX").
			Return(bankAccountRepo, nil)
		repoMockTransaction.EXPECT().
			Create(gomock.Any()).
			Times(1).
			Return(1, nil)
		repoMockBankAccount.EXPECT().
			Update(bankaccountrepo.BankAccount{
				ID:               1,
				OrganizationName: "ACME Corp",
				BalanceCents:     7,
				Iban:             "FR10474608000002006107XXXXX",
				Bic:              "OIVUSCLQXXX",
			}).
			Times(1).
			Return(nil)

		svc := New(repoMockTransaction, repoMockBankAccount, logMock)
		err := svc.BulkTransfer(bulkTransfer)

		assert.Nil(t, err)
	})

	t.Run("Test BulkTransfer return error when user not found", func(t *testing.T) {
		repoMockBankAccount.EXPECT().
			ReadByIban(gomock.Any()).
			Return(bankaccountrepo.BankAccount{}, errors.New("error"))

		svc := New(repoMockTransaction, repoMockBankAccount, logMock)
		err := svc.BulkTransfer(bulkTransfer)

		assert.Error(t, err)
	})

	t.Run("Test BulkTransfer return error when funds are not enough", func(t *testing.T) {

		bankAccountRepoLowBudget := bankaccountrepo.BankAccount{
			ID:               1,
			OrganizationName: "ACME Corp",
			BalanceCents:     1450,
			Iban:             "FR10474608000002006107XXXXX",
			Bic:              "OIVUSCLQXXX",
		}

		repoMockBankAccount.EXPECT().
			ReadByIban(gomock.Any()).
			Return(bankAccountRepoLowBudget, nil)

		svc := New(repoMockTransaction, repoMockBankAccount, logMock)
		err := svc.BulkTransfer(bulkTransfer)

		assert.Error(t, err)
	})
}

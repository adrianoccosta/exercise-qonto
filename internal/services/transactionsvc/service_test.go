package transactionsvc

import (
	"errors"
	"github.com/adrianoccosta/exercise-qonto/internal/domain"
	"github.com/adrianoccosta/exercise-qonto/test/mocks/log"
	"github.com/adrianoccosta/exercise-qonto/test/mocks/repository"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransactionService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := mockrepository.NewMockTransactionRepository(ctrl)
	logMock := mocklog.NewMockLogger(ctrl)

	transactionList := domain.TransactionList{
		{
			Name:             "ACME Corp",
			Iban:             "FR10474608000002006107XXXXX",
			Bic:              "OIVUSCLQXXX",
			CounterPartyName: "Bip Bip",
			CounterPartyIban: "EE383680981021245685",
			CounterPartyBic:  "CRLYFRPPTOU",
			Amount:           14.50,
			Currency:         "EUR",
			Description:      "Wonderland/4410",
		},
	}

	t.Run("Test ReadByFilter return success", func(t *testing.T) {
		repoMock.EXPECT().
			ReadByFilter(gomock.Any()).
			Return(transactionList, nil)

		svc := New(repoMock, logMock)
		res, err := svc.ReadByFilter(make(map[string]string))

		assert.Nil(t, err)
		assert.Equal(t, transactionList, res)
	})

	t.Run("Test ReadByFilter return error", func(t *testing.T) {
		repoMock.EXPECT().
			ReadByFilter(gomock.Any()).
			Return(domain.TransactionList{}, errors.New("error"))

		svc := New(repoMock, logMock)
		_, err := svc.ReadByFilter(make(map[string]string))

		assert.Error(t, err)
	})
}

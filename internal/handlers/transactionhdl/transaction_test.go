package transactionhdl

import (
	"encoding/json"
	"errors"
	"github.com/adrianoccosta/exercise-qonto/internal/domain"
	"github.com/adrianoccosta/exercise-qonto/test/mocks/log"
	"github.com/adrianoccosta/exercise-qonto/test/mocks/services"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTransactionHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	serviceMock := mockservice.NewMockTransactionService(ctrl)
	logMock := mocklog.NewMockLogger(ctrl)

	t.Run("Test read return success", func(t *testing.T) {

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

		filters := map[string]string{
			"name": "ACME Corp",
			"iban": "FR10474608000002006107XXXXX",
		}

		serviceMock.EXPECT().
			ReadByFilter(filters).
			Return(transactionList, nil).Times(1)

		h := New(serviceMock, logMock)
		rr := httptest.NewRecorder()
		r := mux.NewRouter()
		h.Handlers(r)

		req, err := http.NewRequest("GET", "/transaction?name=ACME Corp&iban=FR10474608000002006107XXXXX", nil)
		if err != nil {
			t.Fatal(err)
		}

		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		var actual domain.TransactionList
		err = json.NewDecoder(rr.Body).Decode(&actual)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, transactionList, actual)
	})

	t.Run("Test read return error", func(t *testing.T) {
		filters := map[string]string{}

		serviceMock.EXPECT().ReadByFilter(filters).Return(domain.TransactionList{}, errors.New("error"))
		logMock.EXPECT().WithError(gomock.Any()).Return(logMock).Times(1)
		logMock.EXPECT().Error("error reading transactions").Times(1)

		h := New(serviceMock, logMock)
		rr := httptest.NewRecorder()
		r := mux.NewRouter()
		h.Handlers(r)

		req, err := http.NewRequest("GET", "/transaction", nil)
		if err != nil {
			t.Fatal(err)
		}

		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.NotEmpty(t, rr.Body.String())
	})
}

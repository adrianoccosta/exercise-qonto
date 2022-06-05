package transferhdl

import (
	"errors"
	"github.com/adrianoccosta/exercise-qonto/internal/domain"
	"github.com/adrianoccosta/exercise-qonto/test/mocks/log"
	"github.com/adrianoccosta/exercise-qonto/test/mocks/services"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTransferHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	serviceMock := mockservice.NewMockTransferService(ctrl)
	logMock := mocklog.NewMockLogger(ctrl)

	t.Run("Test transfer return success", func(t *testing.T) {

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

		serviceMock.EXPECT().
			BulkTransfer(bulkTransfer).
			Return(nil).Times(1)

		h := New(serviceMock, logMock)
		rr := httptest.NewRecorder()
		r := mux.NewRouter()
		h.Handlers(r)

		req, err := http.NewRequest("POST", "/transfer/bulk", strings.NewReader("{\"organization_name\": \"ACME Corp\", \"organization_bic\": \"OIVUSCLQXXX\", \"organization_iban\": \"FR10474608000002006107XXXXX\", \"credit_transfers\": [ { \"amount\": \"14.53\", \"currency\": \"EUR\", \"counterparty_name\": \"Bip Bip\", \"counterparty_bic\": \"CRLYFRPPTOU\", \"counterparty_iban\": \"EE383680981021245685\", \"description\": \"Wonderland/4410\"}]}"))
		if err != nil {
			t.Fatal(err)
		}

		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Equal(t, "", rr.Body.String())
	})

	t.Run("Test transfer return error when body is wrong", func(t *testing.T) {

		serviceMock.EXPECT().BulkTransfer(gomock.Any()).Times(0)
		logMock.EXPECT().WithError(gomock.Any()).Return(logMock).Times(1)
		logMock.EXPECT().Error("error parsing message body").Times(1)

		h := New(serviceMock, logMock)
		rr := httptest.NewRecorder()
		r := mux.NewRouter()
		h.Handlers(r)

		req, err := http.NewRequest("POST", "/transfer/bulk", strings.NewReader("{ \"test\":"))
		if err != nil {
			t.Fatal(err)
		}

		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
		assert.NotEmpty(t, rr.Body.String())
	})

	t.Run("Test transfer return error when missing mandatory fields", func(t *testing.T) {

		serviceMock.EXPECT().BulkTransfer(gomock.Any()).Times(0)
		logMock.EXPECT().WithError(gomock.Any()).Return(logMock).Times(1)
		logMock.EXPECT().Error("Missing mandatory fields").Times(1)

		h := New(serviceMock, logMock)
		rr := httptest.NewRecorder()
		r := mux.NewRouter()
		h.Handlers(r)

		req, err := http.NewRequest("POST", "/transfer/bulk", strings.NewReader("{ \"test\": \"ACME Corp\"}"))
		if err != nil {
			t.Fatal(err)
		}

		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
		assert.NotEmpty(t, rr.Body.String())
	})

	t.Run("Test transfer return error", func(t *testing.T) {

		serviceMock.EXPECT().BulkTransfer(gomock.Any()).Return(errors.New("error")).Times(1)
		logMock.EXPECT().WithError(gomock.Any()).Return(logMock).Times(1)
		logMock.EXPECT().Error("error registering bulk transfer").Times(1)

		h := New(serviceMock, logMock)
		rr := httptest.NewRecorder()
		r := mux.NewRouter()
		h.Handlers(r)

		req, err := http.NewRequest("POST", "/transfer/bulk", strings.NewReader("{\"organization_name\": \"ACME Corp\", \"organization_bic\": \"OIVUSCLQXXX\", \"organization_iban\": \"FR10474608000002006107XXXXX\", \"credit_transfers\": [ { \"amount\": \"14.53\", \"currency\": \"EUR\", \"counterparty_name\": \"Bip Bip\", \"counterparty_bic\": \"CRLYFRPPTOU\", \"counterparty_iban\": \"EE383680981021245685\", \"description\": \"Wonderland/4410\"}]}"))
		if err != nil {
			t.Fatal(err)
		}

		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
		assert.NotEmpty(t, rr.Body.String())
	})
}

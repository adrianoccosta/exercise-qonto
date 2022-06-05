package bankaccounthdl

import (
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

func TestBankAccountHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	serviceMock := mockservice.NewMockBankAccountService(ctrl)
	logMock := mocklog.NewMockLogger(ctrl)

	t.Run("Test Create return success", func(t *testing.T) {

		bankAccount := domain.BankAccount{
			Name:    "ACME Corp",
			Balance: 12.40,
			Iban:    "FR10474608000002006107XXXXX",
			Bic:     "OIVUSCLQXXX",
		}

		serviceMock.EXPECT().
			Create(bankAccount).
			Return(nil).Times(1)

		h := New(serviceMock, logMock)
		rr := httptest.NewRecorder()
		r := mux.NewRouter()
		h.Handlers(r)

		req, err := http.NewRequest("POST", "/bank-account", strings.NewReader("{ \"name\": \"ACME Corp\", \"balance\": \"12.40\", \"iban\": \"FR10474608000002006107XXXXX\", \"bic\": \"OIVUSCLQXXX\"}"))
		if err != nil {
			t.Fatal(err)
		}

		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Equal(t, "", rr.Body.String())
	})

	t.Run("Test Create return error when missing mandatory fields", func(t *testing.T) {

		serviceMock.EXPECT().Create(gomock.Any()).Times(0)
		logMock.EXPECT().WithError(gomock.Any()).Return(logMock).Times(1)
		logMock.EXPECT().Error("Missing mandatory fields").Times(1)

		h := New(serviceMock, logMock)
		rr := httptest.NewRecorder()
		r := mux.NewRouter()
		h.Handlers(r)

		req, err := http.NewRequest("POST", "/bank-account", strings.NewReader("{ \"test\": \"ACME Corp\"}"))
		if err != nil {
			t.Fatal(err)
		}

		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.NotEmpty(t, rr.Body.String())
	})
}

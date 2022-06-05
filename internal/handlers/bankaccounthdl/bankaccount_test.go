package bankaccounthdl

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
	"strings"
	"testing"
)

func TestBankAccountHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	serviceMock := mockservice.NewMockBankAccountService(ctrl)
	logMock := mocklog.NewMockLogger(ctrl)

	t.Run("Test create return success", func(t *testing.T) {

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

	t.Run("Test create return error when missing mandatory fields", func(t *testing.T) {

		serviceMock.EXPECT().Create(gomock.Any()).Times(0)
		logMock.EXPECT().WithError(gomock.Any()).Return(logMock).Times(1)
		logMock.EXPECT().Error("error parsing message body").Times(1)

		h := New(serviceMock, logMock)
		rr := httptest.NewRecorder()
		r := mux.NewRouter()
		h.Handlers(r)

		req, err := http.NewRequest("POST", "/bank-account", strings.NewReader("{ \"test\":"))
		if err != nil {
			t.Fatal(err)
		}

		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.NotEmpty(t, rr.Body.String())
	})

	t.Run("Test create return error when missing mandatory fields", func(t *testing.T) {

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

	t.Run("Test create return error", func(t *testing.T) {

		serviceMock.EXPECT().Create(gomock.Any()).Return(errors.New("error")).Times(1)
		logMock.EXPECT().WithError(gomock.Any()).Return(logMock).Times(1)
		logMock.EXPECT().Error("error creating bank account").Times(1)

		h := New(serviceMock, logMock)
		rr := httptest.NewRecorder()
		r := mux.NewRouter()
		h.Handlers(r)

		req, err := http.NewRequest("POST", "/bank-account", strings.NewReader("{ \"name\": \"ACME Corp\", \"balance\": \"12.40\", \"iban\": \"FR10474608000002006107XXXXX\", \"bic\": \"OIVUSCLQXXX\"}"))
		if err != nil {
			t.Fatal(err)
		}

		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.NotEmpty(t, rr.Body.String())
	})

	t.Run("Test read return success", func(t *testing.T) {

		bankAccount := domain.BankAccount{
			Name:    "ACME Corp",
			Balance: 12.40,
			Iban:    "FR10474608000002006107XXXXX",
			Bic:     "OIVUSCLQXXX",
		}

		serviceMock.EXPECT().
			Read("FR10474608000002006107XXXXX").
			Return(bankAccount, nil).Times(1)

		h := New(serviceMock, logMock)
		rr := httptest.NewRecorder()
		r := mux.NewRouter()
		h.Handlers(r)

		req, err := http.NewRequest("GET", "/bank-account/iban/FR10474608000002006107XXXXX", nil)
		if err != nil {
			t.Fatal(err)
		}

		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		var actual domain.BankAccount
		err = json.NewDecoder(rr.Body).Decode(&actual)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, bankAccount, actual)
	})

	t.Run("Test read return error", func(t *testing.T) {

		serviceMock.EXPECT().Read(gomock.Any()).Return(domain.BankAccount{}, errors.New("error"))
		logMock.EXPECT().WithError(gomock.Any()).Return(logMock).Times(1)
		logMock.EXPECT().Error("error reading bank account with iban FR10474608000002006107XXXXX").Times(1)

		h := New(serviceMock, logMock)
		rr := httptest.NewRecorder()
		r := mux.NewRouter()
		h.Handlers(r)

		req, err := http.NewRequest("GET", "/bank-account/iban/FR10474608000002006107XXXXX", nil)
		if err != nil {
			t.Fatal(err)
		}

		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.NotEmpty(t, rr.Body.String())
	})

	t.Run("Test update return success", func(t *testing.T) {

		bankAccount := domain.BankAccount{
			Name:    "ACME Corp",
			Balance: 12.40,
			Iban:    "FR10474608000002006107XXXXX",
			Bic:     "OIVUSCLQXXX",
		}

		serviceMock.EXPECT().
			Update(bankAccount).
			Return(nil).Times(1)

		h := New(serviceMock, logMock)
		rr := httptest.NewRecorder()
		r := mux.NewRouter()
		h.Handlers(r)

		req, err := http.NewRequest("PUT", "/bank-account", strings.NewReader("{ \"name\": \"ACME Corp\", \"balance\": \"12.40\", \"iban\": \"FR10474608000002006107XXXXX\", \"bic\": \"OIVUSCLQXXX\"}"))
		if err != nil {
			t.Fatal(err)
		}

		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Equal(t, "", rr.Body.String())
	})

	t.Run("Test update return error when missing mandatory fields", func(t *testing.T) {

		serviceMock.EXPECT().Update(gomock.Any()).Times(0)
		logMock.EXPECT().WithError(gomock.Any()).Return(logMock).Times(1)
		logMock.EXPECT().Error("Missing mandatory fields").Times(1)

		h := New(serviceMock, logMock)
		rr := httptest.NewRecorder()
		r := mux.NewRouter()
		h.Handlers(r)

		req, err := http.NewRequest("PUT", "/bank-account", strings.NewReader("{ \"test\": \"ACME Corp\"}"))
		if err != nil {
			t.Fatal(err)
		}

		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.NotEmpty(t, rr.Body.String())
	})

	t.Run("Test update return error", func(t *testing.T) {

		serviceMock.EXPECT().Update(gomock.Any()).Return(errors.New("error")).Times(1)
		logMock.EXPECT().WithError(gomock.Any()).Return(logMock).Times(1)
		logMock.EXPECT().Error("error updating bank account").Times(1)

		h := New(serviceMock, logMock)
		rr := httptest.NewRecorder()
		r := mux.NewRouter()
		h.Handlers(r)

		req, err := http.NewRequest("PUT", "/bank-account", strings.NewReader("{ \"name\": \"ACME Corp\", \"balance\": \"12.40\", \"iban\": \"FR10474608000002006107XXXXX\", \"bic\": \"OIVUSCLQXXX\"}"))
		if err != nil {
			t.Fatal(err)
		}

		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.NotEmpty(t, rr.Body.String())
	})

	t.Run("Test delete return success", func(t *testing.T) {

		serviceMock.EXPECT().
			Delete("FR10474608000002006107XXXXX").
			Return(nil).Times(1)

		h := New(serviceMock, logMock)
		rr := httptest.NewRecorder()
		r := mux.NewRouter()
		h.Handlers(r)

		req, err := http.NewRequest("DELETE", "/bank-account/iban/FR10474608000002006107XXXXX", nil)
		if err != nil {
			t.Fatal(err)
		}

		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Test delete return error", func(t *testing.T) {

		serviceMock.EXPECT().Delete(gomock.Any()).Return(errors.New("error"))
		logMock.EXPECT().WithError(gomock.Any()).Return(logMock).Times(1)
		logMock.EXPECT().Error("error deleting bank account").Times(1)

		h := New(serviceMock, logMock)
		rr := httptest.NewRecorder()
		r := mux.NewRouter()
		h.Handlers(r)

		req, err := http.NewRequest("DELETE", "/bank-account/iban/FR10474608000002006107XXXXX", nil)
		if err != nil {
			t.Fatal(err)
		}

		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.NotEmpty(t, rr.Body.String())
	})

}

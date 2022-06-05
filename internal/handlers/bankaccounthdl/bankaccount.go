package bankaccounthdl

import (
	"encoding/json"
	"fmt"
	"github.com/adrianoccosta/exercise-qonto/internal/domain"
	"github.com/adrianoccosta/exercise-qonto/internal/services/bankaccountsvc"
	"github.com/adrianoccosta/exercise-qonto/log"
	"github.com/adrianoccosta/exercise-qonto/tools"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	pathSelection     = "/bank-account"
	pathSelectionIban = "/bank-account/iban/{iban}"
)

// Handler defines the handler interface
type Handler interface {
	Handlers(r *mux.Router)
}

// New returns an implementation of the back account handler
func New(bankAccountService bankaccountsvc.BankAccountService, logger log.Logger) Handler {
	return handler{
		logger:             logger,
		bankAccountService: bankAccountService,
	}
}

type handler struct {
	logger             log.Logger
	bankAccountService bankaccountsvc.BankAccountService
}

func (h handler) Handlers(r *mux.Router) {
	// handlers
	r.HandleFunc(pathSelection, h.create).Methods(http.MethodPost)
	r.HandleFunc(pathSelectionIban, h.read).Methods(http.MethodGet)
	r.HandleFunc(pathSelection, h.update).Methods(http.MethodPut)
	r.HandleFunc(pathSelectionIban, h.delete).Methods(http.MethodDelete)
}

// @Summary create a new bank account if it doesn't exist
// @ID create-bank-account
// @Tags bank account
// @Produce json
// @Param data body domain.BankAccount true "bank account data"
// @Success 201 {string}  string
// @Failure 400 {string}  string
// @Router /v1/bank-account [post]
func (h handler) create(w http.ResponseWriter, r *http.Request) {

	var bankAccount domain.BankAccount
	err := json.NewDecoder(r.Body).Decode(&bankAccount)
	if err != nil {
		h.logger.WithError(err).Error("error parsing message body")
		tools.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err = bankAccount.Validate(); err != nil {
		h.logger.WithError(err).Error("Missing mandatory fields")
		tools.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = h.bankAccountService.Create(bankAccount)
	if err != nil {
		h.logger.WithError(err).Error("error creating bank account")
		tools.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// @Summary read a bank account based on given iban
// @ID read-bank-account-by-iban
// @Tags bank account
// @Produce json
// @Param iban path string true "User iban"
// @Success 200 {object} domain.BankAccount
// @Failure 400 {string}  string
// @Failure 500 {string}  string
// @Router /v1/bank-account/iban/{iban} [get]
func (h handler) read(w http.ResponseWriter, r *http.Request) {

	iban := mux.Vars(r)["iban"]

	if iban == "" {
		h.logger.Error("iban not provided")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bankAccount, err := h.bankAccountService.Read(iban)

	if err != nil {
		h.logger.WithError(err).Error(fmt.Sprintf("error reading banck account with iban %s", iban))
		tools.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	tools.WriteJSON(w, http.StatusOK, bankAccount)
}

// @Summary update an existent bank account
// @ID update-bank-account
// @Tags bank account
// @Produce json
// @Param data body domain.BankAccount true "bank account data"
// @Success 201 {string}  string
// @Failure 400 {string}  string
// @Failure 500 {string}  string
// @Router /v1/bank-account [put]
func (h handler) update(w http.ResponseWriter, r *http.Request) {

	var bankAccount domain.BankAccount
	err := json.NewDecoder(r.Body).Decode(&bankAccount)
	if err != nil {
		h.logger.WithError(err).Error("error parsing message body")
		tools.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = h.bankAccountService.Update(bankAccount)
	if err != nil {
		h.logger.WithError(err).Error("error updating bank account")
		tools.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusCreated)

}

// @Summary delete a bank account based on given iban
// @ID delete-bank-account-by-iban
// @Tags bank account
// @Produce json
// @Param iban path string true "User iban"
// @Success 200 {string}  string
// @Failure 400 {string}  string
// @Failure 500 {string}  string
// @Router /v1/bank-account/iban/{iban} [delete]
func (h handler) delete(w http.ResponseWriter, r *http.Request) {

	iban := mux.Vars(r)["iban"]

	if iban == "" {
		h.logger.Error("iban not provided")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.bankAccountService.Delete(iban)
	if err != nil {
		h.logger.WithError(err).Error("error deleting bank account")
		tools.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

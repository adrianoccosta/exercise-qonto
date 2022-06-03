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

func (h handler) create(w http.ResponseWriter, r *http.Request) {

	var bankAccount domain.BankAccount
	err := json.NewDecoder(r.Body).Decode(&bankAccount)
	if err != nil {
		h.logger.WithError(err).Error("error parsing message body")
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

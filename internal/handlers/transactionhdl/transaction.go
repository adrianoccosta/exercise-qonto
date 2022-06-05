package transactionhdl

import (
	"github.com/adrianoccosta/exercise-qonto/internal/services/transactionsvc"
	"github.com/adrianoccosta/exercise-qonto/log"
	"github.com/adrianoccosta/exercise-qonto/tools"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	pathSelection = "/transaction"
)

// Handler defines the handler interface
type Handler interface {
	Handlers(r *mux.Router)
}

// New returns an implementation of the transaction handler
func New(transactionService transactionsvc.TransactionService, logger log.Logger) Handler {
	return handler{
		logger:             logger,
		transactionService: transactionService,
	}
}

type handler struct {
	logger             log.Logger
	transactionService transactionsvc.TransactionService
}

func (h handler) Handlers(r *mux.Router) {
	// handlers
	r.HandleFunc(pathSelection, h.read).Methods(http.MethodGet)
}

// @Summary Retrieves a bank account based on given iban
// @Description Get details of all transactions
// @Tags transactions
// @ID read-transactions
// @Produce json
// @Param name query string false "transaction search by name"
// @Param iban query string false "transaction search by iban"
// @Param counterparty_name query string false "transaction search by counterparty_name"
// @Success 200 {array} domain.BankAccount
// @Failure 500 {string}  string
// @Router /v1/transaction [get]
func (h handler) read(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()

	filters := make(map[string]string)
	for k, v := range queries {
		if len(v) > 0 {
			filters[k] = v[0]
		}
	}

	transactionList, err := h.transactionService.ReadByFilter(filters)

	if err != nil {
		h.logger.WithError(err).Error("error reading transactions")
		tools.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	tools.WriteJSON(w, http.StatusOK, transactionList)
}

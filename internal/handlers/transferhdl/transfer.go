package transferhdl

import (
	"encoding/json"
	"github.com/adrianoccosta/exercise-qonto/internal/domain"
	"github.com/adrianoccosta/exercise-qonto/internal/services/transfersvc"
	"github.com/adrianoccosta/exercise-qonto/log"
	"github.com/adrianoccosta/exercise-qonto/tools"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	pathSelection = "/transfer/bulk"
)

// Handler defines the handler interface
type Handler interface {
	Handlers(r *mux.Router)
}

// New returns an implementation of the back account handler
func New(transferService transfersvc.TransferService, logger log.Logger) Handler {
	return handler{
		logger:          logger,
		transferService: transferService,
	}
}

type handler struct {
	logger          log.Logger
	transferService transfersvc.TransferService
}

func (h handler) Handlers(r *mux.Router) {
	// handlers
	r.HandleFunc(pathSelection, h.transfer).Methods(http.MethodPost)
}

func (h handler) transfer(w http.ResponseWriter, r *http.Request) {

	var bulkTransfer domain.BulkTransfer

	if err := json.NewDecoder(r.Body).Decode(&bulkTransfer); err != nil {
		h.logger.WithError(err).Error("error parsing message body")
		tools.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}

	if err := h.transferService.BulkTransfer(bulkTransfer); err != nil {
		h.logger.WithError(err).Error("error registering bulk transfer")
		tools.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

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

// @Summary transfer funds in bulk
// @ID create-bulk-transfers
// @Tags transfer
// @Produce json
// @Param data body domain.BulkTransfer true "bulk transfer data"
// @Success 201 {string}  string
// @Failure 422 {string}  string
// @Router /v1/transfer/bulk [post]
func (h handler) transfer(w http.ResponseWriter, r *http.Request) {

	var bulkTransfer domain.BulkTransfer

	if err := json.NewDecoder(r.Body).Decode(&bulkTransfer); err != nil {
		h.logger.WithError(err).Error("error parsing message body")
		tools.WriteError(w, http.StatusUnprocessableEntity, err)
		return
	}

	if err := bulkTransfer.Validate(); err != nil {
		h.logger.WithError(err).Error("Missing mandatory fields")
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

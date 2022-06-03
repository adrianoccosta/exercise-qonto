package transactionsvc

import (
	"github.com/adrianoccosta/exercise-qonto/internal/domain"
	"github.com/adrianoccosta/exercise-qonto/internal/repository/transactionrepo"
	"github.com/adrianoccosta/exercise-qonto/log"
)

// TransactionService Interface for the transaction services
type TransactionService interface {
	ReadByFilter(filters map[string]string) (domain.TransactionList, error)
}

// New returns an instance of the transaction services
func New(transactionrepo transactionrepo.TransactionRepository, logger log.Logger) TransactionService {
	return service{
		logger:          logger,
		transactionrepo: transactionrepo,
	}
}

type service struct {
	logger          log.Logger
	transactionrepo transactionrepo.TransactionRepository
}

// ReadByFilter list of transfers
func (s service) ReadByFilter(filters map[string]string) (domain.TransactionList, error) {
	return s.transactionrepo.ReadByFilter(filters)
}

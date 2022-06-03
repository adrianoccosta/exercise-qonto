package transfersvc

import (
	"errors"
	"github.com/adrianoccosta/exercise-qonto/internal/domain"
	"github.com/adrianoccosta/exercise-qonto/internal/repository/bankaccountrepo"
	"github.com/adrianoccosta/exercise-qonto/internal/repository/transactionrepo"
	"github.com/adrianoccosta/exercise-qonto/log"
)

// TransferService Interface for the transfer services
type TransferService interface {
	BulkTransfer(data domain.BulkTransfer) error
}

// New returns an instance of the back account services
func New(transactionrepo transactionrepo.TransactionRepository, bankAccountRepo bankaccountrepo.BankAccountRepository, logger log.Logger) TransferService {
	return service{
		logger:          logger,
		transactionrepo: transactionrepo,
		bankAccountRepo: bankAccountRepo,
	}
}

type service struct {
	logger          log.Logger
	transactionrepo transactionrepo.TransactionRepository
	bankAccountRepo bankaccountrepo.BankAccountRepository
}

func (s service) BulkTransfer(data domain.BulkTransfer) error {

	bankAccount, err := s.bankAccountRepo.ReadByIban(data.OrganizationIban)

	if err != nil {
		return err
	}

	var totalAmount float64 = 0
	for _, creditTransfer := range data.CreditTransfers {
		totalAmount += creditTransfer.Amount
	}

	if bankAccount.BalanceCents < int(totalAmount*100) {
		return errors.New("Insufficient credits to complete the transfer")
	}

	bankAccount.BalanceCents -= int(totalAmount * 100)

	go s.registerTransfers(bankAccount, data)

	return nil
}

func (s service) registerTransfers(bankAccount bankaccountrepo.BankAccount, data domain.BulkTransfer) {
	for _, creditTransfer := range data.CreditTransfers {
		s.transactionrepo.Create(transactionrepo.Transaction{
			CounterPartyName: creditTransfer.CounterPartyName,
			CounterPartyIban: creditTransfer.CounterPartyIban,
			CounterPartyBic:  creditTransfer.CounterPartyBic,
			AmountCents:      int(creditTransfer.Amount * 100),
			AmountCurrency:   creditTransfer.Currency,
			BankAccountID:    bankAccount.ID,
			Description:      creditTransfer.Description,
		})
	}

	s.bankAccountRepo.Update(bankAccount)
}

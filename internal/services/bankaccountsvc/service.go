package bankaccountsvc

import (
	"github.com/adrianoccosta/exercise-qonto/internal/domain"
	"github.com/adrianoccosta/exercise-qonto/internal/repository/bankaccountrepo"
	"github.com/adrianoccosta/exercise-qonto/log"
)

// BankAccountService Interface for the back account services
type BankAccountService interface {
	Create(data domain.BankAccount) error
	Read(iban string) (domain.BankAccount, error)
	Update(data domain.BankAccount) error
	Delete(iban string) error
}

// New returns an instance of the back account services
func New(bankAccountRepo bankaccountrepo.BankAccountRepository, logger log.Logger) BankAccountService {
	return service{
		logger:          logger,
		bankAccountRepo: bankAccountRepo,
	}
}

type service struct {
	logger          log.Logger
	bankAccountRepo bankaccountrepo.BankAccountRepository
}

// Create new bank account
func (s service) Create(data domain.BankAccount) error {
	_, err := s.bankAccountRepo.Create(bankaccountrepo.BankAccount{
		OrganizationName: data.Name,
		BalanceCents:     int(data.Balance * 100),
		Iban:             data.Iban,
		Bic:              data.Bic,
	})

	return err
}

// Read a bank account
func (s service) Read(iban string) (domain.BankAccount, error) {
	info, err := s.bankAccountRepo.ReadByIban(iban)

	if err != nil {
		return domain.BankAccount{}, err
	}

	return domain.BankAccount{
		Name:    info.OrganizationName,
		Balance: float64(info.BalanceCents) / 100,
		Iban:    info.Iban,
		Bic:     info.Bic,
	}, nil
}

// Update the values of a bank account
func (s service) Update(data domain.BankAccount) error {
	info, err := s.bankAccountRepo.ReadByIban(data.Iban)

	if err != nil {
		return err
	}

	info.OrganizationName = data.Name
	info.BalanceCents = int(data.Balance * 100)
	info.Bic = data.Bic

	return s.bankAccountRepo.Update(info)
}

// Delete a bank account
func (s service) Delete(iban string) error {
	return s.bankAccountRepo.DeleteByIban(iban)
}

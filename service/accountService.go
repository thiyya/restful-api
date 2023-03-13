package service

import (
	"fmt"
	"net/http"
	"restful-api/global"

	"restful-api/model"
	"restful-api/repo"
)

type AccountService struct {
	accountRepo repo.AccountRepo
}

func NewAccountService(options ...func(*AccountService)) *AccountService {
	as := &AccountService{}
	for _, o := range options {
		o(as)
	}
	return as
}

func WithAccountRepo(accountRepo repo.AccountRepo) func(*AccountService) {
	return func(s *AccountService) {
		s.accountRepo = accountRepo
	}
}

func (a *AccountService) Transfer(debitAccountId, creditAccountId string, amount float64) error {
	debit, err := a.GetAccountById(debitAccountId)
	if err != nil {
		err.(*global.Error).SetErrorMessage(fmt.Sprintf("debit %s", err.Error()))
		return err
	}
	err = a.ValidateTransfer(debit, amount)
	if err != nil {
		return err
	}
	credit, err := a.GetAccountById(creditAccountId)
	if err != nil {
		err.(*global.Error).SetErrorMessage(fmt.Sprintf("credit %s", err.Error()))
		return err
	}
	debit.Balance -= amount
	a.accountRepo.Update(debit)

	credit.Balance += amount
	a.accountRepo.Update(credit)
	return nil
}

func (a *AccountService) ValidateTransfer(debit *model.Account, amount float64) error {
	if debit.Balance < amount {
		return global.NewError(http.StatusForbidden, global.InsufficientFunds, "insufficient funds")
	}

	return nil
}

func (a *AccountService) GetAllAccounts() map[string]*model.Account {
	return a.accountRepo.GetAllAccounts()
}

func (a *AccountService) GetAccountById(accountId string) (*model.Account, error) {
	return a.accountRepo.GetAccountById(accountId)
}

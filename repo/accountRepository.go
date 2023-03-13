package repo

import (
	"net/http"
	"restful-api/global"
	"sync"

	"restful-api/model"
)

type AccountWriteRepo interface {
	Update(account *model.Account)
}

type AccountReadRepo interface {
	GetAllAccounts() map[string]*model.Account
	GetAccountById(accountId string) (*model.Account, error)
}

type AccountRepo interface {
	AccountWriteRepo
	AccountReadRepo
}

var _ AccountRepo = (*accountRepository)(nil)

type accountRepository struct {
	accountMap map[string]*model.Account
	mutex      sync.RWMutex
}

func NewAccountRepository(options ...func(*accountRepository)) *accountRepository {
	ar := &accountRepository{}
	for _, o := range options {
		o(ar)
	}
	return ar
}

func WithAccounts(accountMap map[string]*model.Account) func(*accountRepository) {
	return func(s *accountRepository) {
		s.accountMap = accountMap
	}
}

func (a *accountRepository) Update(account *model.Account) {
	a.mutex.Lock()
	a.accountMap[account.ID] = account
	a.mutex.Unlock()
}

func (a *accountRepository) GetAllAccounts() map[string]*model.Account {
	a.mutex.RLock()
	accounts := a.accountMap
	a.mutex.RUnlock()
	return accounts
}

func (a *accountRepository) GetAccountById(accountId string) (*model.Account, error) {
	a.mutex.RLock()
	account, ok := a.accountMap[accountId]
	a.mutex.RUnlock()
	if ok {
		return account, nil
	}

	return nil, global.NewError(http.StatusNotFound, global.NotFoundErr, "account not found")
}

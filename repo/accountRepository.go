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
	sync.RWMutex
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
	a.Lock()
	defer a.Unlock()

	a.accountMap[account.ID] = account
}

func (a *accountRepository) GetAllAccounts() map[string]*model.Account {
	a.RLock()
	defer a.RUnlock()

	result := make(map[string]*model.Account, len(a.accountMap))
	for k, v := range a.accountMap {
		result[k] = v
	}
	return result
}

func (a *accountRepository) GetAccountById(accountId string) (*model.Account, error) {
	a.RLock()
	defer a.RUnlock()

	account, ok := a.accountMap[accountId]
	if ok {
		return account, nil
	}

	return nil, global.NewError(http.StatusNotFound, global.NotFoundErr, "account not found")
}

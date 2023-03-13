package handler

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"restful-api/dto"
	"restful-api/global"
	"restful-api/service"
)

type AccountHandler struct {
	accountService *service.AccountService
}

func NewAccountHandler(options ...func(*AccountHandler)) *AccountHandler {
	as := &AccountHandler{}
	for _, o := range options {
		o(as)
	}
	return as
}

func WithAccountService(accountService *service.AccountService) func(*AccountHandler) {
	return func(s *AccountHandler) {
		s.accountService = accountService
	}
}

func (a *AccountHandler) Transfer() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		input := dto.TransferRequest{}
		json.NewDecoder(r.Body).Decode(&input)
		err := validator.New().Struct(input)
		if err != nil {
			log.Printf("err occurred while parsing transfer input: %s \n", err.Error())
			err = global.NewError(http.StatusBadRequest, global.InvalidParams, err.Error())
			err.(*global.Error).WriteError(w)
			return
		}

		err = a.accountService.Transfer(input.DebitAccountId, input.CreditAccountId, input.Amount)
		if err != nil {
			log.Printf("err occurred while transfer : %s \n", err.Error())
			err.(*global.Error).WriteError(w)

			return
		}

		writeResponse(w, dto.TransferResponse{Success: true}, http.StatusCreated)
	})
}

func (a *AccountHandler) Accounts() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := dto.AccountsResponse{}
		for _, account := range a.accountService.GetAllAccounts() {
			res.Accounts = append(res.Accounts, dto.Account{
				ID:      account.ID,
				Name:    account.Name,
				Balance: account.Balance,
			})
		}
		writeResponse(w, res, http.StatusOK)
	})
}

func (a *AccountHandler) GetAccountById() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		account, err := a.accountService.GetAccountById(id)
		if err != nil {
			log.Printf("err occurred while getting account for id %s : %s \n", id, err.Error())
			err.(*global.Error).WriteError(w)

			return
		}
		writeResponse(w, dto.Account{
			ID:      account.ID,
			Name:    account.Name,
			Balance: account.Balance,
		}, http.StatusOK)
	})
}

func writeResponse(w http.ResponseWriter, v any, responseCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseCode)
	b, _ := json.Marshal(v)
	w.Write(b)
}

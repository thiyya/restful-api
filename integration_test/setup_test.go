package integration_test

import (
	"github.com/gorilla/mux"
	"net/http"
	"restful-api/handler"
	"restful-api/model"
	"restful-api/repo"
	"restful-api/service"
)

func Prepare() {
	accountMap := map[string]*model.Account{
		"1": {
			ID:      "1",
			Name:    "Low Balance",
			Balance: 1,
		},
		"2": {
			ID:      "2",
			Name:    "Name2",
			Balance: 1000,
		},
		"3": {
			ID:      "3",
			Name:    "Name3",
			Balance: 2000,
		},
		"4": {
			ID:      "4",
			Name:    "Name4",
			Balance: 4000,
		},
	}
	accHandler := handler.NewAccountHandler(handler.WithAccountService(service.NewAccountService(service.WithAccountRepo(repo.NewAccountRepository(repo.WithAccounts(accountMap))))))
	mux := mux.NewRouter()
	mux.Handle("/transfer", accHandler.Transfer()).Methods("POST")
	mux.Handle("/account", accHandler.Accounts()).Methods("GET")
	mux.Handle("/account/{id}", accHandler.GetAccountById()).Methods("GET")
	srv := &http.Server{Addr: ":8080", Handler: mux}
	go srv.ListenAndServe()
}

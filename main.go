package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"restful-api/handler"
	"restful-api/model"
	"restful-api/repo"
	"restful-api/service"
)

func main() {
	accRepo := repo.NewAccountRepository(repo.WithAccounts(consumeAccounts()))
	log.Println("accounts consumed successfully")
	accServ := service.NewAccountService(service.WithAccountRepo(accRepo))
	accHandler := handler.NewAccountHandler(handler.WithAccountService(accServ))

	mux := mux.NewRouter()
	mux.Handle("/transfer", accHandler.Transfer()).Methods("POST")
	mux.Handle("/account", accHandler.Accounts()).Methods("GET")
	mux.Handle("/account/{id}", accHandler.GetAccountById()).Methods("GET")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()
	log.Println("server started successfully, ready to transfer")

	stopC := make(chan os.Signal)
	signal.Notify(stopC, os.Interrupt)
	<-stopC

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	log.Println("server stopping ...")
	defer cancel()

	log.Fatal(srv.Shutdown(ctx))
}

func consumeAccounts() map[string]*model.Account {
	response, err := http.Get("https://git.io/Jm76h")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	defer response.Body.Close()
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var accountList []model.Account
	err = json.Unmarshal(responseData, &accountList)
	if err != nil {
		log.Fatal(err)
	}
	accountMap := make(map[string]*model.Account)
	for i := 0; i < len(accountList); i++ {
		accountMap[accountList[i].ID] = &accountList[i]
	}

	return accountMap
}

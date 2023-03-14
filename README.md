## restful-api

It is an api that provides money transfer between debit and credit accounts. 
* Accounts are fetched from https://git.io/Jm76h and ingested into the local memory
* After merging, the server starts and notifies you if it is ready to transfer
* Transfer requests are validated to confirm both debit and credit account, and that the transfer will not result in a negative balance

## How to run:
```bash
go run main.go
```

## How to test:
```bash
go test ./...
```

** For convenience, I have shared a postman collection under the docs directory
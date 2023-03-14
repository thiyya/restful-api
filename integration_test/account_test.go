package integration_test

import (
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"io"
	"net/http"
	"restful-api/dto"
	"restful-api/global"
	"strings"
	"testing"
)

const (
	testCase1 = `{
"debitAccountId":"2",
"creditAccountId":"3",
"amount":100
}`
	testCase2 = `{
"debitAccountId":"1",
"creditAccountId":"4",
"amount":100
}`
	testCase3 = `{
"debitAccountId":"1",
"creditAccountId":"4",
"amount":0
}`
	testCase4 = `{
"debitAccountId":"2",
"creditAccountId":"1234",
"amount":100
}`
	testCase5 = `{
"debitAccountId":"1321",
"creditAccountId":"4",
"amount":100
}`
)

func Test_Handler(t *testing.T) {
	Prepare()
	Convey("Transfer", t, func(c C) {
		Convey("Transfer Successfully", func(c C) {
			req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/transfer", strings.NewReader(testCase1))
			c.So(err, ShouldBeNil)
			res, err := http.DefaultClient.Do(req)
			c.So(err, ShouldBeNil)
			c.So(res.StatusCode, ShouldEqual, http.StatusOK)
			resBody, err := io.ReadAll(res.Body)
			c.So(err, ShouldBeNil)
			tr := dto.TransferResponse{}
			err = json.Unmarshal(resBody, &tr)
			c.So(err, ShouldBeNil)
			c.So(tr.Success, ShouldEqual, true)
			req, err = http.NewRequest(http.MethodGet, "http://localhost:8080/account/3", nil)
			c.So(err, ShouldBeNil)
			res, err = http.DefaultClient.Do(req)
			c.So(err, ShouldBeNil)
			c.So(res.StatusCode, ShouldEqual, http.StatusOK)
			resBody, err = io.ReadAll(res.Body)
			c.So(err, ShouldBeNil)
			a := dto.Account{}
			err = json.Unmarshal(resBody, &a)
			c.So(err, ShouldBeNil)
			c.So(a.ID, ShouldEqual, "3")
			c.So(a.Name, ShouldEqual, "Name3")
			c.So(a.Balance, ShouldEqual, 2100)
		})
		Convey("Transfer Insufficient Funds", func(c C) {
			req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/transfer", strings.NewReader(testCase2))
			c.So(err, ShouldBeNil)
			res, err := http.DefaultClient.Do(req)
			c.So(err, ShouldBeNil)
			c.So(res.StatusCode, ShouldEqual, http.StatusForbidden)
			resBody, err := io.ReadAll(res.Body)
			c.So(err, ShouldBeNil)
			e := dto.Error{}
			err = json.Unmarshal(resBody, &e)
			c.So(err, ShouldBeNil)
			c.So(e.Item, ShouldEqual, global.InsufficientFunds)
			c.So(e.Message, ShouldEqual, "insufficient funds")

			req, err = http.NewRequest(http.MethodGet, "http://localhost:8080/account/4", nil)
			c.So(err, ShouldBeNil)
			res, err = http.DefaultClient.Do(req)
			c.So(err, ShouldBeNil)
			c.So(res.StatusCode, ShouldEqual, http.StatusOK)
			resBody, err = io.ReadAll(res.Body)
			c.So(err, ShouldBeNil)
			a := dto.Account{}
			err = json.Unmarshal(resBody, &a)
			c.So(err, ShouldBeNil)
			c.So(a.ID, ShouldEqual, "4")
			c.So(a.Name, ShouldEqual, "Name4")
			c.So(a.Balance, ShouldEqual, 4000)
		})
		Convey("Transfer Invalid Params", func(c C) {
			req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/transfer", strings.NewReader(testCase3))
			c.So(err, ShouldBeNil)
			res, err := http.DefaultClient.Do(req)
			c.So(err, ShouldBeNil)
			c.So(res.StatusCode, ShouldEqual, http.StatusBadRequest)
			resBody, err := io.ReadAll(res.Body)
			c.So(err, ShouldBeNil)
			e := dto.Error{}
			err = json.Unmarshal(resBody, &e)
			c.So(err, ShouldBeNil)
			c.So(e.Item, ShouldEqual, global.InvalidParams)
			c.So(e.Message, ShouldContainSubstring, "Error:Field validation")
		})
		Convey("Transfer Not Found Credit Account", func(c C) {
			req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/transfer", strings.NewReader(testCase4))
			c.So(err, ShouldBeNil)
			res, err := http.DefaultClient.Do(req)
			c.So(err, ShouldBeNil)
			c.So(res.StatusCode, ShouldEqual, http.StatusNotFound)
			resBody, err := io.ReadAll(res.Body)
			c.So(err, ShouldBeNil)
			e := dto.Error{}
			err = json.Unmarshal(resBody, &e)
			c.So(err, ShouldBeNil)
			c.So(e.Item, ShouldEqual, global.NotFoundErr)
			c.So(e.Message, ShouldEqual, "credit account not found")
		})
		Convey("Transfer Not Found Debit Account", func(c C) {
			req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/transfer", strings.NewReader(testCase5))
			c.So(err, ShouldBeNil)
			res, err := http.DefaultClient.Do(req)
			c.So(err, ShouldBeNil)
			c.So(res.StatusCode, ShouldEqual, http.StatusNotFound)
			resBody, err := io.ReadAll(res.Body)
			c.So(err, ShouldBeNil)
			e := dto.Error{}
			err = json.Unmarshal(resBody, &e)
			c.So(err, ShouldBeNil)
			c.So(e.Item, ShouldEqual, global.NotFoundErr)
			c.So(e.Message, ShouldEqual, "debit account not found")
		})
	})
	Convey("Get accounts", t, func(c C) {
		Convey("Get All Accounts", func(c C) {
			req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/account", nil)
			c.So(err, ShouldBeNil)
			res, err := http.DefaultClient.Do(req)
			c.So(err, ShouldBeNil)
			c.So(res.StatusCode, ShouldEqual, http.StatusOK)
			resBody, err := io.ReadAll(res.Body)
			c.So(err, ShouldBeNil)
			ar := dto.AccountsResponse{}
			err = json.Unmarshal(resBody, &ar)
			c.So(err, ShouldBeNil)
			c.So(len(ar.Accounts), ShouldEqual, 4)
		})
		Convey("Get Account By Id", func(c C) {
			req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/account/1", nil)
			c.So(err, ShouldBeNil)
			res, err := http.DefaultClient.Do(req)
			c.So(err, ShouldBeNil)
			c.So(res.StatusCode, ShouldEqual, http.StatusOK)
			resBody, err := io.ReadAll(res.Body)
			c.So(err, ShouldBeNil)
			a := dto.Account{}
			err = json.Unmarshal(resBody, &a)
			c.So(err, ShouldBeNil)
			c.So(a.ID, ShouldEqual, "1")
			c.So(a.Name, ShouldEqual, "Low Balance")
			c.So(a.Balance, ShouldEqual, 1)
		})
	})
}

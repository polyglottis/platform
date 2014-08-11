package test

import (
	"bytes"
	"testing"

	"github.com/polyglottis/platform/language"
	"github.com/polyglottis/platform/user"
)

type Tester struct {
	server user.Server
	*testing.T
}

func NewTester(server user.Server, t *testing.T) *Tester {
	return &Tester{
		server: server,
		T:      t,
	}
}

var Account = &user.NewAccountRequest{
	Name:         "testUser",
	MainLanguage: language.English.Code,
	Email:        "test@test.com",
	PasswordHash: []byte("testPW"),
}

func (t *Tester) All() {
	t.NotExist(Account.Name)

	a := t.NewTwice(Account)

	t.Equals(a, Account)

	t.Get(Account.Name, a)
}

func (t *Tester) NotExist(n user.Name) {
	_, err := t.server.GetAccount(n)
	if err != user.AccountNotFound {
		t.Errorf("Test DB should be empty, but got error '%v'", err)
	}
}

func (t *Tester) NewTwice(r *user.NewAccountRequest) *user.Account {
	a, err := t.server.NewAccount(r)
	if err != nil {
		t.Fatal(err)
	}
	if a == nil {
		t.Fatal("Test account should have been created an not be nil")
	}
	_, err = t.server.NewAccount(r)
	if err == nil {
		t.Errorf("Creating the same account twice should trigger an error")
	}
	return a
}

func (t *Tester) Equals(a *user.Account, r *user.NewAccountRequest) {
	if a.Name != r.Name || a.MainLanguage != r.MainLanguage || a.Email != r.Email || !bytes.Equal(a.PasswordHash, r.PasswordHash) {
		t.Errorf("These accounts should coincide, except for the Active field: %+v != %+v", a, r)
	}
}

func (t *Tester) Get(n user.Name, check *user.Account) {
	a, err := t.server.GetAccount(n)
	if err != nil {
		t.Fatal(err)
	}
	if a == nil {
		t.Fatal("An account should get returned here")
	}
	if !check.Equals(a) {
		t.Errorf("These accounts should coincide: %+v != %+v", a, check)
	}
}

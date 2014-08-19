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

	t.GetByEmail(Account.Email, a)

	a.Email = "newTest@test.com"
	a.Active = false
	a.MainLanguage = language.Unknown.Code
	a.PasswordHash = []byte("updatedPW")
	t.UpdateAccount(a)

	t.GetNotFound("other")
	t.GetByEmailNotFound("other@email.com")

	t.Tokens(Account.Name)
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
	t.Check(check, a)
}

func (t *Tester) Check(check, other *user.Account) {
	if !check.Equals(other) {
		t.Errorf("These accounts should coincide: %+v != %+v", check, other)
	}
}

func (t *Tester) GetByEmail(email string, check *user.Account) {
	a, err := t.server.GetAccountByEmail(email)
	if err != nil {
		t.Fatal(err)
	}
	if a == nil {
		t.Fatal("An account should get returned here")
	}
	t.Check(check, a)
}

func (t *Tester) GetNotFound(n user.Name) {
	_, err := t.server.GetAccount(n)
	if err != user.AccountNotFound {
		t.Fatal("Account should not be found: %v", n)
	}
}

func (t *Tester) GetByEmailNotFound(email string) {
	_, err := t.server.GetAccountByEmail(email)
	if err != user.AccountNotFound {
		t.Fatal("Account should not be found: %s", email)
	}
}

func (t *Tester) UpdateAccount(a *user.Account) {
	err := t.server.UpdateAccount(a)
	if err != nil {
		t.Fatal(err)
	}
	t.Get(a.Name, a)
}

func (t *Tester) Tokens(n user.Name) {
	token, err := t.server.NewToken(n)
	if err != nil {
		t.Fatal(err)
	}

	valid, err := t.server.ValidToken(n, token)
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Error("New tokens should always be valid")
	}

	err = t.server.DeleteToken(n, token)
	if err != nil {
		t.Fatal(err)
	}

	valid, err = t.server.ValidToken(n, token)
	if err != nil {
		t.Fatal(err)
	}
	if valid {
		t.Error("Deleted tokens should not be valid")
	}
}

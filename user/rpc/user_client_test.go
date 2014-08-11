package rpc

import (
	"fmt"
	"testing"

	"github.com/polyglottis/platform/user"
	"github.com/polyglottis/platform/user/test"
)

var addr = ":1234"

// server is just a proxy for testing.
type server struct {
	accounts map[user.Name]*user.Account
}

func (s *server) NewAccount(r *user.NewAccountRequest) (*user.Account, error) {
	if _, exists := s.accounts[r.Name]; exists {
		return nil, fmt.Errorf("User name %s is already taken", r.Name)
	}

	a := user.NewAccount(r)
	s.accounts[r.Name] = a
	return a, nil
}

func (s *server) GetAccount(n user.Name) (*user.Account, error) {
	if a, ok := s.accounts[n]; ok {
		return a, nil
	}
	return nil, fmt.Errorf("Account not found")
}

func TestServerAndClient(t *testing.T) {

	testServer := NewUserServer(&server{
		accounts: make(map[user.Name]*user.Account),
	}, addr)

	err := testServer.RegisterAndListen()
	if err != nil {
		t.Fatal(err)
	}

	go testServer.Accept()

	client, err := NewClient(addr)
	if err != nil {
		t.Fatal(err)
	}

	tester := test.NewTester(client, t)
	tester.All()
}

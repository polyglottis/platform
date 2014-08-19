package rpc

import (
	"fmt"
	"testing"

	"github.com/polyglottis/platform/user"
	"github.com/polyglottis/platform/user/test"
	"github.com/polyglottis/rand"
)

var addr = ":1234"

// server is just a proxy for testing.
type server struct {
	accounts map[user.Name]*user.Account
	tokens   map[user.Name][]string
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

func (s *server) GetAccountByEmail(email string) (*user.Account, error) {
	for _, a := range s.accounts {
		if a.Email == email {
			return a, nil
		}
	}
	return nil, fmt.Errorf("Account not found")
}

func (s *server) NewToken(n user.Name) (string, error) {
	token, err := rand.Id(12)
	if err != nil {
		return "", err
	}
	s.tokens[n] = append(s.tokens[n], token)
	return token, nil
}

func (s *server) ValidToken(n user.Name, token string) (bool, error) {
	if tokens, ok := s.tokens[n]; ok {
		for _, t := range tokens {
			if t == token {
				return true, nil
			}
		}
	}
	return false, nil
}

func (s *server) DeleteToken(n user.Name, token string) error {
	if tokens, ok := s.tokens[n]; ok {
		for i, t := range tokens {
			if t == token {
				s.tokens[n] = append(tokens[:i], tokens[i+1:]...)
				break
			}
		}
	}
	return nil
}

func TestServerAndClient(t *testing.T) {

	testServer := NewUserServer(&server{
		accounts: make(map[user.Name]*user.Account),
		tokens:   make(map[user.Name][]string),
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

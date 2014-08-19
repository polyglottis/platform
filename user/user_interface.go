package user

import (
	"errors"
)

// Server is the interface a user server should comply to.
type Server interface {
	NewAccount(r *NewAccountRequest) (*Account, error)
	GetAccount(n Name) (*Account, error)
	GetAccountByEmail(email string) (*Account, error)
	UpdateAccount(a *Account) error

	NewToken(n Name) (string, error)
	ValidToken(n Name, token string) (bool, error)
	DeleteToken(n Name, token string) error
}

var AccountNotFound = errors.New("Account not found")

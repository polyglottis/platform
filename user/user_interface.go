package user

import (
	"errors"
)

// Server is the interface a user server should comply to.
type Server interface {
	NewAccount(r *NewAccountRequest) (*Account, error)
	GetAccount(n Name) (*Account, error)
}

var AccountNotFound = errors.New("Account not found")

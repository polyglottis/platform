// Package user contains the type definitions for user accounts and the user server interface.
package user

import (
	"github.com/polyglottis/platform/language"
)

type Name string

type Account struct {
	Name         Name
	MainLanguage language.Code
	Email        string
	Active       bool
	PasswordHash []byte
}

type NewAccountRequest struct {
	Name         Name
	MainLanguage language.Code
	Email        string
	PasswordHash []byte
}

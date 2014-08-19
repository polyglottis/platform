// Package user contains the type definitions for user accounts and the user server interface.
package user

import (
	"regexp"
	
	"github.com/polyglottis/platform/i18n"
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

// ValidPassword checks whether a password is strong enough.
// If not, a message is returned (a non-empty i18n key).
func ValidPassword(passwd string) (bool, i18n.Key) {
	if len(passwd) < 8 {
		return false, i18n.Key("Password too short")
	}
	return true, ""
}

var validName = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]*$`)

// ValidName checks whether a username is acceptable.
// If not, a message is returned (a non-empty i18n key).
func ValidName(name string) (bool, i18n.Key) {
	if len(name) < 3 {
		return false, i18n.Key("Username too short")
	}
	if validName.MatchString(name) {
		return true, ""
	} else {
		return false, "Username must start with a letter and consist entirely of letters, numbers and underscores."
	}
}
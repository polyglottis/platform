package user

import (
	"bytes"
)

func NewAccount(r *NewAccountRequest) *Account {
	return &Account{
		Name:         r.Name,
		UILocale:     r.UILocale,
		Email:        r.Email,
		Active:       true,
		PasswordHash: r.PasswordHash,
	}
}

func (a *Account) Equals(b *Account) bool {
	if a == nil {
		return b == nil
	} else if b == nil {
		return false
	}
	// now a != nil and b != nil
	if a.Name != b.Name ||
		a.UILocale != b.UILocale ||
		a.MainLanguage != b.MainLanguage ||
		a.Email != b.Email ||
		a.Active != b.Active ||
		!bytes.Equal(a.PasswordHash, b.PasswordHash) {
		return false
	}
	return true
}

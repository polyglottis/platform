// Package password provides password hashing and password checks.
package password

import (
	"code.google.com/p/go.crypto/bcrypt"
	
	"github.com/polyglottis/platform/user"
)

const bcrypt_cost = 12 // between 4 and 31 (both inclusive)

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt_cost)
}

// Check returns true if the password is correct.
func Check(password string, a *user.Account) bool {
	err := bcrypt.CompareHashAndPassword(a.PasswordHash, []byte(password))
	// (err == nil) if and only if the password matches
	return err == nil
}

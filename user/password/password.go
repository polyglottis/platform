package user

import (
	"code.google.com/p/go.crypto/bcrypt"
)

const bcrypt_cost = 12 // between 4 and 31 (both inclusive)

func getPasswordHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt_cost)
}

// PasswordCorrect returns true if the password is correct.
func (a *Account) PasswordCorrect(password string) bool {
	err := bcrypt.CompareHashAndPassword(a.passwordHash, []byte(password))
	// (err == nil) if and only if the password matches
	return err == nil
}

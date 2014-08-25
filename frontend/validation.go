package frontend

import (
	"crypto/rand"
	"math/big"
	mathRand "math/rand"
	"regexp"
	"time"

	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

var emailRegex = regexp.MustCompile(`(?i)^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,4}$`)

func validEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func sleep() {
	var i int64
	var limit int64 = 1000
	bigInt, err := rand.Int(rand.Reader, big.NewInt(limit))
	if err == nil {
		i = bigInt.Int64()
	} else {
		i = mathRand.Int63n(limit)
	}
	t := 1*time.Second + time.Duration(i)*time.Millisecond
	time.Sleep(t)
}

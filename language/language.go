// Package language defines language codes, the language struct, and the language server interface.
package language

import (
	"errors"
)

type Code string

type Language struct {
	Code      Code
	ISO_639_1 string
	ISO_639_3 string
	WikiData  string
}

var Unknown = Language{Code: "unknown"}
var None = Language{Code: "none"}
var English = Language{
	Code:      "en",
	ISO_639_1: "en",
	ISO_639_3: "eng",
	WikiData:  "",
}

// Server is the interface a language server should comply to.
type Server interface {
	List() ([]Code, error)
	GetCode(maybeCode string) (Code, error)
}

var CodeNotFound = errors.New("Language code not found")

type CodeSlice []Code

func (s CodeSlice) Len() int           { return len(s) }
func (s CodeSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s CodeSlice) Less(i, j int) bool { return s[i] < s[j] }

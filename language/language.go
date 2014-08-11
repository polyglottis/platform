// Package language defines language codes, the language struct, and the language server interface.
package language

type Code string

type Language struct {
	Code      Code
	ISO_693_1 string
	ISO_693_3 string
	ISO_693_6 string
}

var Unknown = Language{Code: "unknown"}
var None = Language{Code: "none"}
var English = Language{
	Code:      "en",
	ISO_693_1: "en",
	ISO_693_3: "eng",
	ISO_693_6: "eng",
}

// Server is the interface a language server should comply to.
type Server interface {
	GetCode(code string) (Code, error)
}

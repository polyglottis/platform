package test

import (
	"fmt"

	"github.com/polyglottis/platform/language"
)

func English(server language.Server) error {
	code, err := server.GetCode("en")
	if err != nil {
		return err
	}
	if code != language.English.Code {
		return fmt.Errorf("GetCode should return English language code when asked for 'en'")
	}
	return nil
}

func Invalid(server language.Server) error {
	invalid := "this is not a valid code"
	code, err := server.GetCode(invalid)
	if err == nil {
		return fmt.Errorf("GetCode('%s') should trigger an error", invalid)
	}
	if code != language.Unknown.Code {
		return fmt.Errorf("GetCode('%s') should return language Unknown", invalid)
	}
	return nil
}

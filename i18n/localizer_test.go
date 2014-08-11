package i18n

import (
	"testing"

	"github.com/polyglottis/platform/language"
)

func TestNewLocalizer(t *testing.T) {
	loc := NewLocalizer(language.English.Code)

	testData := []struct {
		In  Key
		Out string
	}{{
		In:  "test",
		Out: "test",
	}}

	for i, test := range testData {
		got := loc.GetText(test.In)
		if got != test.Out {
			t.Errorf("Got %v but expected %v (test %d)", got, test.Out, i)
		}
	}
}

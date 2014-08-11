package i18n

import (
	"github.com/polyglottis/platform/language"
)

type Key string

type Localizer interface {
	GetText(Key) string
}

type localizer struct {
	locale language.Code
}

func NewLocalizer(locale language.Code) Localizer {
	return &localizer{
		locale: locale,
	}
}

func (loc *localizer) GetText(key Key) string {
	return string(key)
}

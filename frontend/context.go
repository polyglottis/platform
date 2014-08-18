package frontend

import (
	"net/http"
	"net/url"

	"github.com/gorilla/mux"

	"github.com/polyglottis/platform/content"
	"github.com/polyglottis/platform/language"
)

type Context struct {
	Locale    language.Code
	Vars      map[string]string
	Query     url.Values
	ExtractId content.ExtractId
	LanguageA language.Code
	LanguageB language.Code
}

func ReadContext(r *http.Request) (*Context, error) {
	return &Context{
		Locale: language.English.Code,
		Vars:   mux.Vars(r),
		Query:  r.URL.Query(),
	}, nil
}

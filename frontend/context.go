package frontend

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/polyglottis/platform/language"
)

type Context struct {
	Locale language.Code
	Vars   map[string]string
}

func ReadContext(r *http.Request) (*Context, error) {
	return &Context{
		Locale: language.English.Code,
		Vars:   mux.Vars(r),
	}, nil
}

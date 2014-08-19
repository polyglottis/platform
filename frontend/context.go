package frontend

import (
	"net/http"
	"net/url"

	"github.com/gorilla/mux"

	"github.com/polyglottis/platform/content"
	"github.com/polyglottis/platform/i18n"
	"github.com/polyglottis/platform/language"
	"github.com/polyglottis/platform/user"
)

type Context struct {
	Locale language.Code
	Vars   map[string]string
	Query  url.Values
	Form   url.Values
	Url string
	User   user.Name

	Defaults url.Values          // default form values
	Errors   map[string]i18n.Key // errors on form submit

	ExtractId content.ExtractId
	LanguageA language.Code
	LanguageB language.Code
}

func ReadContext(r *http.Request, s *Session) (*Context, error) {
	c := &Context{
		Locale: language.English.Code,
		Vars:   mux.Vars(r),
		Query:  r.URL.Query(),
		Url: r.URL.String(),
	}
	
	if u :=  s.GetAccount(); u != nil {
			c.User = u.Name
			}
	
	return c, nil
}

func ReadContextWithForm(r *http.Request, s *Session) (*Context, error) {
	c, err := ReadContext(r, s)
	if err != nil {
		return nil, err
	}

	err = r.ParseForm()
	if err != nil {
		return nil, err
	}

	c.Form = r.PostForm
	return c, nil
}

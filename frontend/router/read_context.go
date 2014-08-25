package router

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/polyglottis/platform/frontend"
	"github.com/polyglottis/platform/frontend/handle"
)

func ReadContext(r *http.Request, s *handle.Session) (*frontend.Context, error) {
	c := &frontend.Context{
		Locale: "en-us",
		Vars:   mux.Vars(r),
		Query:  r.URL.Query(),
		Url:    r.URL.String(),
		Host:   r.Host,
	}

	if r.TLS == nil {
		c.Protocol = "http"
	} else {
		c.Protocol = "https"
	}

	if u := s.GetAccount(); u != nil {
		c.User = u.Name
	}

	return c, nil
}

func ReadContextWithForm(r *http.Request, s *handle.Session) (*frontend.Context, error) {
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

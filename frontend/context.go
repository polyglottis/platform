package frontend

import (
	"net/url"

	"github.com/polyglottis/platform/i18n"
	"github.com/polyglottis/platform/user"
)

type Context struct {
	Locale   string
	Vars     map[string]string
	Query    url.Values
	Form     url.Values
	Url      string
	User     user.Name
	Protocol string
	Host     string

	Email string // for password reset

	Defaults url.Values // default form values
	Errors   ErrorMap   // errors on form submit
}

type ErrorMap map[string]i18n.Key

func (c *Context) ProtocolAndHost() string {
	return c.Protocol + "://" + c.Host
}

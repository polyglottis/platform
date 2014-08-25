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

	Defaults url.Values          // default form values
	Errors   map[string]i18n.Key // errors on form submit
}

func (c *Context) ProtocolAndHost() string {
	return c.Protocol + "://" + c.Host
}

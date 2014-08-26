package frontend

import (
	"net/url"
	"strconv"

	"github.com/polyglottis/platform/content"
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

func (c *Context) LoggedIn() bool {
	return len(c.User) != 0
}

func (c *Context) IsFocusOnA() bool {
	return c.Query.Get("focus") != "b"
}

func (c *Context) FlavorId(key string) (content.FlavorId, error) {
	i, err := strconv.Atoi(c.Query.Get(key))
	if err != nil {
		return 0, err
	}
	return content.FlavorId(i), nil
}

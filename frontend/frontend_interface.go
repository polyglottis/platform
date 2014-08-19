// Package frontend defines the interface a frontend server should comply to.
package frontend

import (
	"github.com/polyglottis/platform/content"
)

// Server is the interface a frontend server should comply to.
type Server interface {
	Home(context *Context) ([]byte, error)
	NotFound(context *Context) ([]byte, error)

	Extract(context *Context, e *content.Extract) ([]byte, error)
	Flavor(context *Context, e *content.Extract, a, b *FlavorTriple) ([]byte, error)

	EditText(context *Context, e *content.Extract, a, b *content.Flavor) ([]byte, error)

	SignUp(context *Context) ([]byte, error)
	SignIn(context *Context) ([]byte, error)
}

type FlavorTriple struct {
	Audio      *content.Flavor
	Text       *content.Flavor
	Transcript *content.Flavor
}

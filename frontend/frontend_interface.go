// Package frontend defines the interface a frontend server should comply to.
package frontend

import (
	"github.com/polyglottis/platform/content"
	"github.com/polyglottis/platform/language"
	"github.com/polyglottis/platform/user"
)

// Server is the interface a frontend server should comply to.
type Server interface {
	SetLanguageList([]language.Code) error

	Home(context *Context) ([]byte, error)
	Error(context *Context) ([]byte, error)
	NotFound(context *Context) ([]byte, error)

	Flavor(context *Context, e *content.Extract, a, b *FlavorTriple) ([]byte, error)

	NewExtract(context *Context) ([]byte, error)
	NewFlavor(context *Context, e *content.Extract) ([]byte, error)
	EditText(context *Context, e *content.Extract, a, b *content.Flavor) ([]byte, error)

	SignUp(context *Context) ([]byte, error)
	SignIn(context *Context) ([]byte, error)
	ForgotPassword(context *Context) ([]byte, error)
	PasswordSent(context *Context) ([]byte, error)
	ResetPassword(context *Context) ([]byte, error)
	PasswordResetEmail(c *Context, a *user.Account, token string) ([]byte, error)
}

type FlavorTriple struct {
	Audio      *content.Flavor
	Text       *content.Flavor
	Transcript *content.Flavor
}

func (t *FlavorTriple) Language() language.Code {
	switch {
	case t.Text != nil:
		return t.Text.Language
	case t.Audio != nil:
		return t.Audio.Language
	case t.Transcript != nil:
		return t.Transcript.Language
	default:
		return language.Unknown.Code
	}
}

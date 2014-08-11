package backend

import (
	"github.com/polyglottis/platform/content"
	"github.com/polyglottis/platform/language"
	"github.com/polyglottis/platform/user"
)

type Configuration struct {
	Content  content.Server
	Language language.Server
	User     user.Server
}

type Engine struct {
	*Configuration
}

func NewEngine(c *Configuration) *Engine {
	return &Engine{
		Configuration: c,
	}
}

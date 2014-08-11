// Package rpc provides the rpc language client used by the Polyglottis Application
// and a simple language server wrapper.
package rpc

import (
	"github.com/polyglottis/platform/language"
	"github.com/polyglottis/rpc"
)

// LanguageServer is the rpc language server object.
// It is used internally by NewLanguageServer, and is not meant to be instantiated directly.
// It needs to be exported for rpc to work.
type LanguageServer struct {
	s language.Server
}

// NewLanguageServer creates an rpc language server, forwarding calls to s, and listening on tcp address addr.
func NewLanguageServer(s language.Server, addr string) *rpc.Server {
	return rpc.NewServer("LanguageServer", &LanguageServer{s}, addr)
}

func (s *LanguageServer) GetCode(code string, reply *language.Code) error {
	newCode, err := s.s.GetCode(code)
	if err != nil {
		return err
	}
	*reply = newCode
	return nil
}

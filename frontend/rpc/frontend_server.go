// Package rpc provides the rpc frontend client used by the Polyglottis Application
// and a simple frontend server wrapper.
package rpc

import (
	"github.com/polyglottis/platform/frontend"
	"github.com/polyglottis/rpc"
)

// FrontendServer is the rpc frontend server object.
// It is used internally by NewFrontendServer, and is not meant to be instantiated directly.
// It needs to be exported for rpc to work.
type FrontendServer struct {
	s frontend.Server
}

// NewFrontendServer creates an rpc frontend server, forwarding calls to s, and listening on tcp address addr.
func NewFrontendServer(s frontend.Server, addr string) *rpc.Server {
	return rpc.NewServer("FrontendServer", &FrontendServer{s}, addr)
}

func (s *FrontendServer) Home(context *frontend.Context, r *[]byte) error {
	bytes, err := s.s.Home(context)
	if err != nil {
		return err
	}
	*r = bytes
	return nil
}

func (s *FrontendServer) NotFound(context *frontend.Context, r *[]byte) (err error) {
	*r, err = s.s.NotFound(context)
	return
}

func (s *FrontendServer) Extract(ce *ContextExtract, r *[]byte) (err error) {
	*r, err = s.s.Extract(ce.Context, ce.Extract)
	return
}

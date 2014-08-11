// Package rpc provides the rpc user client used by the Polyglottis Application
// and a simple user server wrapper.
package rpc

import (
	"github.com/polyglottis/platform/user"
	"github.com/polyglottis/rpc"
)

// UserServer is the rpc user server object.
// It is used internally by NewUserServer, and is not meant to be instantiated directly.
// It needs to be exported for rpc to work.
type UserServer struct {
	s user.Server
}

// NewUserServer creates an rpc user server, forwarding calls to s, and listening on tcp address addr.
func NewUserServer(s user.Server, addr string) *rpc.Server {
	return rpc.NewServer("UserServer", &UserServer{s}, addr)
}

func (s *UserServer) NewAccount(r *user.NewAccountRequest, a *user.Account) error {
	acc, err := s.s.NewAccount(r)
	if err != nil {
		return err
	}
	*a = *acc
	return nil
}

func (s *UserServer) GetAccount(n user.Name, a *user.Account) error {
	acc, err := s.s.GetAccount(n)
	if err != nil {
		return err
	}
	*a = *acc
	return nil
}

// Package rpc provides the rpc frontend client used by the Polyglottis Application
// and a simple frontend server wrapper.
package rpc

import (
	"github.com/polyglottis/platform/content"
	"github.com/polyglottis/platform/frontend"
	"github.com/polyglottis/platform/language"
	"github.com/polyglottis/platform/user"
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

func (s *FrontendServer) SetLanguageList(list []language.Code, nothing *bool) error {
	return s.s.SetLanguageList(list)
}
func (s *FrontendServer) Home(context *frontend.Context, r *[]byte) (err error) {
	*r, err = s.s.Home(context)
	return
}
func (s *FrontendServer) Error(context *frontend.Context, r *[]byte) (err error) {
	*r, err = s.s.Error(context)
	return
}
func (s *FrontendServer) NotFound(context *frontend.Context, r *[]byte) (err error) {
	*r, err = s.s.NotFound(context)
	return
}

func (s *FrontendServer) SignUp(context *frontend.Context, r *[]byte) (err error) {
	*r, err = s.s.SignUp(context)
	return
}
func (s *FrontendServer) SignIn(context *frontend.Context, r *[]byte) (err error) {
	*r, err = s.s.SignIn(context)
	return
}

func (s *FrontendServer) ForgotPassword(context *frontend.Context, r *[]byte) (err error) {
	*r, err = s.s.ForgotPassword(context)
	return
}
func (s *FrontendServer) PasswordSent(context *frontend.Context, r *[]byte) (err error) {
	*r, err = s.s.PasswordSent(context)
	return
}
func (s *FrontendServer) ResetPassword(context *frontend.Context, r *[]byte) (err error) {
	*r, err = s.s.ResetPassword(context)
	return
}

type ContextFlavorTriples struct {
	Context *frontend.Context
	Extract *content.Extract
	A       *frontend.FlavorTriple
	B       *frontend.FlavorTriple
}
type ContextFlavors struct {
	Context *frontend.Context
	Extract *content.Extract
	A       *content.Flavor
	B       *content.Flavor
}
type ContextExtract struct {
	Context *frontend.Context
	Extract *content.Extract
}

func (s *FrontendServer) Flavor(cf *ContextFlavorTriples, r *[]byte) (err error) {
	*r, err = s.s.Flavor(cf.Context, cf.Extract, cf.A, cf.B)
	return
}

func (s *FrontendServer) NewExtract(context *frontend.Context, r *[]byte) (err error) {
	*r, err = s.s.NewExtract(context)
	return
}
func (s *FrontendServer) NewFlavor(ce *ContextExtract, r *[]byte) (err error) {
	*r, err = s.s.NewFlavor(ce.Context, ce.Extract)
	return
}
func (s *FrontendServer) EditText(cf *ContextFlavors, r *[]byte) (err error) {
	*r, err = s.s.EditText(cf.Context, cf.Extract, cf.A, cf.B)
	return
}

type AccountToken struct {
	Context *frontend.Context
	Account *user.Account
	Token   string
}

func (s *FrontendServer) PasswordResetEmail(a *AccountToken, r *[]byte) (err error) {
	*r, err = s.s.PasswordResetEmail(a.Context, a.Account, a.Token)
	return
}

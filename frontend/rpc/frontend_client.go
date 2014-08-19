package rpc

import (
	"github.com/polyglottis/platform/content"
	"github.com/polyglottis/platform/frontend"
	"github.com/polyglottis/platform/user"
	"github.com/polyglottis/rpc"
)

// Client is a frontend.Server which forwards method calls to another frontend.Server over rpc.
type Client struct {
	rpc *rpc.Client
}

func NewClient(addr string) (*Client, error) {
	c, err := rpc.NewClient("Frontend", "tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Client{rpc: c}, nil
}

func (c *Client) call(funcName string, context *frontend.Context) ([]byte, error) {
	var response []byte
	err := c.rpc.Call("FrontendServer."+funcName, context, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (c *Client) callExtract(funcName string, context *frontend.Context, e *content.Extract) ([]byte, error) {
	var response []byte
	err := c.rpc.Call("FrontendServer."+funcName, &ContextExtract{
		Context: context,
		Extract: e,
	}, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (c *Client) Home(context *frontend.Context) ([]byte, error) {
	return c.call("Home", context)
}

func (c *Client) NotFound(context *frontend.Context) ([]byte, error) {
	return c.call("NotFound", context)
}

func (c *Client) Extract(context *frontend.Context, extract *content.Extract) ([]byte, error) {
	return c.callExtract("Extract", context, extract)
}

func (c *Client) Flavor(context *frontend.Context, extract *content.Extract, a, b *frontend.FlavorTriple) ([]byte, error) {
	var response []byte
	err := c.rpc.Call("FrontendServer.Flavor", &ContextFlavorTriples{
		Context: context,
		Extract: extract,
		A:       a,
		B:       b,
	}, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (c *Client) EditText(context *frontend.Context, extract *content.Extract, a, b *content.Flavor) ([]byte, error) {
	var response []byte
	err := c.rpc.Call("FrontendServer.EditText", &ContextFlavors{
		Context: context,
		Extract: extract,
		A:       a,
		B:       b,
	}, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (c *Client) SignUp(context *frontend.Context) ([]byte, error) {
	return c.call("SignUp", context)
}
func (c *Client) SignIn(context *frontend.Context) ([]byte, error) {
	return c.call("SignIn", context)
}

func (c *Client) ForgotPassword(context *frontend.Context) ([]byte, error) {
	return c.call("ForgotPassword", context)
}
func (c *Client) PasswordSent(context *frontend.Context) ([]byte, error) {
	return c.call("PasswordSent", context)
}
func (c *Client) ResetPassword(context *frontend.Context) ([]byte, error) {
	return c.call("ResetPassword", context)
}
func (c *Client) PasswordResetEmail(context *frontend.Context, a *user.Account, token string) (b []byte, err error) {
	err = c.rpc.Call("FrontendServer.PasswordResetEmail", &AccountToken{
		Context: context,
		Account: a,
		Token:   token,
	}, &b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

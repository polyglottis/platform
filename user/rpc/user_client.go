package rpc

import (
	"log"

	"github.com/polyglottis/platform/user"
	"github.com/polyglottis/rpc"
)

// Client is a user.Server which forwards method calls to another user.Server over rpc.
type Client struct {
	rpc *rpc.Client
}

func NewClient(addr string) (*Client, error) {
	c, err := rpc.NewClient("User", "tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Client{rpc: c}, nil
}

func (c *Client) NewAccount(r *user.NewAccountRequest) (*user.Account, error) {
	log.Printf("Creating new account %+v", r)
	a := new(user.Account)
	err := c.rpc.Call("UserServer.NewAccount", r, a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (c *Client) GetAccount(n user.Name) (*user.Account, error) {
	a := new(user.Account)
	err := c.rpc.Call("UserServer.GetAccount", n, a)
	if err != nil {
		if err.Error() == "Account not found" {
			return nil, user.AccountNotFound
		}
		return nil, err
	}
	return a, nil
}

func (c *Client) GetAccountByEmail(email string) (*user.Account, error) {
	a := new(user.Account)
	err := c.rpc.Call("UserServer.GetAccountByEmail", email, a)
	if err != nil {
		if err.Error() == "Account not found" {
			return nil, user.AccountNotFound
		}
		return nil, err
	}
	return a, nil
}
func (c *Client) UpdateAccount(a *user.Account) error {
	return c.rpc.Call("UserServer.UpdateAccount", a, nil)
}

func (c *Client) NewToken(n user.Name) (token string, err error) {
	err = c.rpc.Call("UserServer.NewToken", n, &token)
	return
}

type NamedToken struct {
	Name  user.Name
	Token string
}

func (c *Client) ValidToken(n user.Name, token string) (valid bool, err error) {
	err = c.rpc.Call("UserServer.ValidToken", &NamedToken{Name: n, Token: token}, &valid)
	return
}
func (c *Client) DeleteToken(n user.Name, token string) error {
	return c.rpc.Call("UserServer.DeleteToken", &NamedToken{Name: n, Token: token}, nil)
}

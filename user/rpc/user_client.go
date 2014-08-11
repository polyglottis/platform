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

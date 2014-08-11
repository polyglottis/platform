package rpc

import (
	"github.com/polyglottis/platform/language"
	"github.com/polyglottis/rpc"
)

type Client struct {
	c *rpc.Client
}

func NewClient(addr string) (*Client, error) {
	c, err := rpc.NewClient("Language", "tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Client{c: c}, nil
}

func (c *Client) GetCode(code string) (language.Code, error) {
	reply := new(language.Code)
	err := c.c.Call("LanguageServer.GetCode", code, reply)
	if err != nil {
		return language.Unknown.Code, err
	}
	return *reply, nil
}

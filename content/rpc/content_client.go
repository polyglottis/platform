package rpc

import (
	"github.com/polyglottis/platform/content"
	"github.com/polyglottis/platform/user"
	"github.com/polyglottis/rpc"
)

type Client struct {
	rpc *rpc.Client
}

func NewClient(addr string) (*Client, error) {
	c, err := rpc.NewClient("Content", "tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Client{rpc: c}, nil
}

func (c *Client) NewExtract(author user.Name, e *content.Extract) error {
	var id content.ExtractId
	err := c.rpc.Call("ContentServer.NewExtract", &ExtractRequest{
		Author:  author,
		Extract: e,
	}, &id)
	if err != nil {
		return err
	}

	e.SetId(id)
	e.SetFlavorLanguagesAndTypes()
	return nil
}

func (c *Client) NewFlavor(author user.Name, f *content.Flavor) error {
	var id content.FlavorId
	err := c.rpc.Call("ContentServer.NewFlavor", &FlavorRequest{
		Author: author,
		Flavor: f,
	}, &id)
	if err != nil {
		return err
	}

	f.SetId(id)
	f.SetLanguage(f.Language)
	f.SetType(f.Type)
	return nil
}

func (c *Client) ExtractList() ([]*content.Extract, error) {
	var list []*content.Extract
	err := c.rpc.Call("ContentServer.ExtractList", false, &list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (c *Client) GetExtract(id content.ExtractId) (*content.Extract, error) {
	e := new(content.Extract)
	err := c.rpc.Call("ContentServer.GetExtract", id, e)
	if err != nil {
		if err.Error() == content.ErrNotFound.Error() {
			return nil, content.ErrNotFound
		}
		return nil, err
	}
	return e, nil
}

func (c *Client) GetExtractId(slug string) (content.ExtractId, error) {
	var id content.ExtractId
	err := c.rpc.Call("ContentServer.GetExtractId", slug, &id)
	if err != nil {
		if err.Error() == content.ErrNotFound.Error() {
			return "", content.ErrNotFound
		}
		return "", err
	}
	return id, nil
}

func (c *Client) UpdateExtract(author user.Name, e *content.Extract) error {
	return c.rpc.Call("ContentServer.UpdateExtract", &ExtractRequest{
		Author:  author,
		Extract: e,
	}, nil)
}

func (c *Client) UpdateFlavor(author user.Name, f *content.Flavor) error {
	return c.rpc.Call("ContentServer.UpdateFlavor", &FlavorRequest{
		Author: author,
		Flavor: f,
	}, nil)
}

func (c *Client) InsertOrUpdateUnits(author user.Name, u []*content.Unit) error {
	return c.rpc.Call("ContentServer.InsertOrUpdateUnits", &UnitsRequest{
		Author: author,
		Units:  u,
	}, nil)
}

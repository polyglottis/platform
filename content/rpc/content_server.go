// Package rpc provides the rpc content client used by the Polyglottis Application
// and a simple content server wrapper.
package rpc

import (
	"github.com/polyglottis/platform/content"
	"github.com/polyglottis/platform/language"
	"github.com/polyglottis/platform/user"
	"github.com/polyglottis/rpc"
)

// ContentServer is the rpc content server object.
// It is used internally by NewContentServer, and is not meant to be instantiated directly.
// It needs to be exported for rpc to work.
type ContentServer struct {
	s content.Server
}

// NewContentServer creates an rpc content server, forwarding calls to s, and listening on tcp address addr.
func NewContentServer(s content.Server, addr string) *rpc.Server {
	return rpc.NewServer("ContentServer", &ContentServer{s}, addr)
}

type ExtractRequest struct {
	Author  user.Name
	Extract *content.Extract
}

type FlavorRequest struct {
	Author user.Name
	Flavor *content.Flavor
}

type UnitsRequest struct {
	Author user.Name
	Units  []*content.Unit
}

func (s *ContentServer) NewExtract(r *ExtractRequest, id *content.ExtractId) error {
	err := s.s.NewExtract(r.Author, r.Extract)
	if err != nil {
		return err
	}
	*id = r.Extract.Id
	return nil
}

func (s *ContentServer) NewFlavor(r *FlavorRequest, id *content.FlavorId) error {
	err := s.s.NewFlavor(r.Author, r.Flavor)
	if err != nil {
		return err
	}
	*id = r.Flavor.Id
	return nil
}

func (s *ContentServer) ExtractList(nothing bool, list *[]*content.Extract) (err error) {
	*list, err = s.s.ExtractList()
	return
}

func (s *ContentServer) ExtractLanguages(nothing bool, list *[]language.Code) (err error) {
	*list, err = s.s.ExtractLanguages()
	return
}

func (s *ContentServer) GetExtract(id content.ExtractId, e *content.Extract) error {
	extract, err := s.s.GetExtract(id)
	if err != nil {
		return err
	}
	*e = *extract
	return nil
}

func (s *ContentServer) GetExtractId(slug string, id *content.ExtractId) error {
	newId, err := s.s.GetExtractId(slug)
	if err != nil {
		return err
	}
	*id = newId
	return nil
}

func (s *ContentServer) UpdateExtract(r *ExtractRequest, nothing *bool) error {
	return s.s.UpdateExtract(r.Author, r.Extract)
}

func (s *ContentServer) UpdateFlavor(r *FlavorRequest, nothing *bool) error {
	return s.s.UpdateFlavor(r.Author, r.Flavor)
}

func (s *ContentServer) InsertOrUpdateUnits(r *UnitsRequest, nothing *bool) error {
	return s.s.InsertOrUpdateUnits(r.Author, r.Units)
}

package rpc

import (
	"fmt"
	"sort"
	"testing"

	"github.com/polyglottis/platform/content"
	"github.com/polyglottis/platform/content/test"
	"github.com/polyglottis/platform/user"
	"github.com/polyglottis/rand"
)

var addr = ":1234"

// server is just a proxy for testing.
type server struct {
	extracts map[content.ExtractId]*content.Extract
}

func (s *server) NewExtract(author user.Name, e *content.Extract) error {
	strId, err := rand.Id(5)
	if err != nil {
		return err
	}

	id := content.ExtractId(strId)

	if _, exists := s.extracts[id]; exists {
		return fmt.Errorf("Internal error")
	}

	s.extracts[id] = e
	e.SetId(id)
	return nil
}

func (s *server) NewFlavor(author user.Name, f *content.Flavor) error {
	if e, ok := s.extracts[f.ExtractId]; ok {
		var max content.FlavorId
		for _, flav := range e.Flavors {
			if max < flav.Id {
				max = flav.Id
			}
		}
		e.Flavors = append(e.Flavors, f)
		f.SetId(max + 1)
		return nil
	}
	return content.ErrNotFound
}

func (s *server) ExtractList() ([]*content.Extract, error) {
	list := make([]*content.Extract, 0)
	for _, e := range s.extracts {
		list = append(list, e)
	}
	return list, nil
}

func (s *server) GetExtract(id content.ExtractId) (*content.Extract, error) {
	if e, ok := s.extracts[id]; ok {
		return e, nil
	}
	return nil, content.ErrNotFound
}

func (s *server) GetExtractId(slug string) (content.ExtractId, error) {
	for id, e := range s.extracts {
		if e.UrlSlug == slug {
			return id, nil
		}
	}
	return "", content.ErrNotFound
}

func (s *server) UpdateExtract(author user.Name, e *content.Extract) error {
	if _, ok := s.extracts[e.Id]; ok {
		s.extracts[e.Id] = e
		return nil
	} else {
		return content.ErrNotFound
	}
}

func (s *server) UpdateFlavor(author user.Name, f *content.Flavor) error {
	if _, ok := s.extracts[f.ExtractId]; ok {
		flavors := s.extracts[f.ExtractId].Flavors
		for i, flav := range flavors {
			if f.Id == flav.Id {
				flavors[i] = f
				return nil
			}
		}
	}
	return content.ErrNotFound
}

func (s *server) InsertOrUpdateUnits(author user.Name, units []*content.Unit) error {
	for _, u := range units {
		err := s.InsertOrUpdateUnit(author, u)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *server) InsertOrUpdateUnit(author user.Name, u *content.Unit) error {
	if u.BlockId < 1 || u.Id < 1 || (u.BlockId == 1 && u.Id != 1) {
		return content.ErrInvalidInput
	}
	if _, ok := s.extracts[u.ExtractId]; ok {
		flavors := s.extracts[u.ExtractId].Flavors
		for _, flav := range flavors {
			if u.FlavorId == flav.Id {
				for i, block := range flav.Blocks {
					if block[0].BlockId == u.BlockId {
						for j, unit := range block {
							if unit.Id == u.Id {
								block[j] = u
								return nil
							}
						}
						flav.Blocks[i] = append(flav.Blocks[i], u)
						sort.Sort(content.UnitSlice(flav.Blocks[i]))
						return nil
					}
				}
				flav.Blocks = append(flav.Blocks, []*content.Unit{u})
				sort.Sort(content.BlockSlice(flav.Blocks))
				return nil
			}
		}
	}
	return content.ErrNotFound
}

func TestServerAndClient(t *testing.T) {

	testServer := NewContentServer(&server{
		extracts: make(map[content.ExtractId]*content.Extract),
	}, addr)

	err := testServer.RegisterAndListen()
	if err != nil {
		t.Fatal(err)
	}

	go testServer.Accept()

	client, err := NewClient(addr)
	if err != nil {
		t.Fatal(err)
	}

	tester := test.NewTester(client, t)
	tester.All()
}

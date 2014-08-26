package rpc

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/polyglottis/platform/content"
	"github.com/polyglottis/platform/content/test"
	"github.com/polyglottis/platform/language"
	"github.com/polyglottis/platform/user"
	"github.com/polyglottis/rand"
)

var addr = ":1234"

// server is just a proxy for testing.
type server struct {
	extracts map[content.ExtractId]*content.Extract
}

func (s *server) NewExtract(author user.Name, e *content.Extract) error {
	if valid, _ := content.ValidSlug(e.UrlSlug); !valid || !content.ValidExtractType(e.Type) {
		return content.ErrInvalidInput
	}

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
	e.SetFlavorLanguagesAndTypes()
	return nil
}

func (s *server) NewFlavor(author user.Name, f *content.Flavor) error {
	if len(f.Type) == 0 || len(f.Language) == 0 {
		return content.ErrInvalidInput
	}
	if e, ok := s.extracts[f.ExtractId]; ok {
		if _, ok := e.Flavors[f.Language]; !ok {
			e.Flavors[f.Language] = make(content.FlavorByType)
		}
		fByType := e.Flavors[f.Language]
		flavors := fByType[f.Type]
		var max content.FlavorId
		for _, flav := range flavors {
			if max < flav.Id {
				max = flav.Id
			}
		}
		fByType[f.Type] = append(flavors, f)
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

func (s *server) ExtractLanguages() ([]language.Code, error) {
	list := make([]language.Code, 0)
	for _, e := range s.extracts {
		for lang := range e.Flavors {
			list = append(list, lang)
		}
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

func (s *server) ExtractsMatching(q *content.Query) ([]content.ExtractId, error) {
	matching := make([]content.ExtractId, 0)
	for id, e := range s.extracts {
		if q.LanguageA != "" {
			if _, ok := e.Flavors[q.LanguageA]; !ok {
				continue
			}
		}
		if q.LanguageB != "" {
			if _, ok := e.Flavors[q.LanguageB]; !ok {
				continue
			}
		}
		if q.ExtractType != "" {
			if e.Type != q.ExtractType {
				continue
			}
		}
		matching = append(matching, id)
	}
	return matching, nil
}

func (s *server) UpdateExtract(author user.Name, e *content.Extract) error {
	if _, ok := s.extracts[e.Id]; ok {
		s.extracts[e.Id].Type = e.Type
		s.extracts[e.Id].Metadata = e.Metadata
		return nil
	} else {
		return content.ErrNotFound
	}
}

func (s *server) UpdateFlavor(author user.Name, f *content.Flavor) error {
	if _, ok := s.extracts[f.ExtractId]; ok {
		if fByType, ok := s.extracts[f.ExtractId].Flavors[f.Language]; ok {
			if flavors, ok := fByType[f.Type]; ok {
				for i, flav := range flavors {
					if f.Id == flav.Id {
						flavors[i] = f
						return nil
					}
				}
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
	if len(u.Language) == 0 || len(u.FlavorType) == 0 {
		return content.ErrInvalidInput
	}
	if _, ok := s.extracts[u.ExtractId]; ok {
		if fByType, ok := s.extracts[u.ExtractId].Flavors[u.Language]; ok {
			if flavors, ok := fByType[u.FlavorType]; ok {
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
		}
	}
	return content.ErrNotFound
}

func TestClientImplementsInterface(t *testing.T) {
	client, err := NewClient(addr)
	if err != nil {
		t.Fatal(err)
	}

	cType := reflect.TypeOf(client)
	typeOfServer := reflect.TypeOf((*content.Server)(nil)).Elem()
	if !cType.Implements(typeOfServer) {
		t.Fatal("RPC client does not implement content.Server")
	}
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

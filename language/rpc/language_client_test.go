package rpc

import (
	"testing"

	"github.com/polyglottis/platform/language"
	"github.com/polyglottis/platform/language/test"
)

type server struct{}

func (s *server) GetCode(code string) (language.Code, error) {
	if code == "en" {
		return language.English.Code, nil
	} else {
		return "", language.CodeNotFound
	}
}

func (s *server) List() ([]language.Code, error) {
	return []language.Code{language.English.Code}, nil
}

func TestServerAndClient(t *testing.T) {
	addr := ":1234"

	testServer := NewLanguageServer(&server{}, addr)

	err := testServer.RegisterAndListen()
	if err != nil {
		t.Fatal(err)
	}

	go testServer.Accept()

	c, err := NewClient(addr)
	if err != nil {
		t.Fatal(err)
	}

	test.All(c, t)
}

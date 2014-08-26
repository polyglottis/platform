// Package api contains the definition of all web services to extract data from the Polyglottis Platform.
package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/polyglottis/platform/backend"
	"github.com/polyglottis/platform/content"
	"github.com/polyglottis/platform/server"
)

type Server struct {
	*backend.Engine
}

func NewServer(engine *backend.Engine) *Server {
	return &Server{
		Engine: engine,
	}
}

func (s *Server) RegisterServices(r *mux.Router) {
	r.HandleFunc("/", errorHandler(s.Root))
	r.HandleFunc("/extract/by-id/{id}", errorHandler(s.ExtractById)).Methods("GET")
	r.HandleFunc("/extract/list", errorHandler(s.ExtractList)).Methods("GET")
	r.HandleFunc("/extract/languages", errorHandler(s.ExtractLanguages)).Methods("GET")
	r.HandleFunc("/extract/search", errorHandler(s.ExtractSearch)).Methods("GET")
}

func errorHandler(f func(*http.Request) ([]byte, error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				server.Recovered(r, w, rec)
			}
		}()
		bytes, err := f(r)
		if err == nil {
			w.Write(bytes)
		} else {
			Error(r, w, err)
		}
	}
}

func call(f func(io.Writer) error) ([]byte, error) {
	buffer := new(bytes.Buffer)
	err := f(buffer)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (s *Server) Root(r *http.Request) ([]byte, error) {
	return call(func(w io.Writer) error {
		w.Write([]byte("hello"))
		return nil
	})
}

func (s *Server) ExtractById(r *http.Request) ([]byte, error) {
	return call(func(w io.Writer) error {
		vars := mux.Vars(r)
		id := content.ExtractId(vars["id"])
		if len(id) == 0 {
			return fmt.Errorf("id parameter is mandatory")
		}
		e, err := s.Content.GetExtract(id)
		if err != nil {
			return err
		}
		return json.NewEncoder(w).Encode(e)
	})
}

func (s *Server) ExtractList(r *http.Request) ([]byte, error) {
	return call(func(w io.Writer) error {
		list, err := s.Content.ExtractList()
		if err != nil {
			return err
		}
		return json.NewEncoder(w).Encode(list)
	})
}

func Error(r *http.Request, w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	log.Printf("Error: %v (while serving %v)", err, r)
	w.Write([]byte(err.Error()))
}

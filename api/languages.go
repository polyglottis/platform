// Package api contains the definition of all web services to extract data from the Polyglottis Platform.
package api

import (
	"encoding/json"
	"io"
	"net/http"
)

func (s *Server) ExtractLanguages(r *http.Request) ([]byte, error) {
	return call(func(w io.Writer) error {
		list, err := s.Content.ExtractLanguages()
		if err != nil {
			return err
		}
		return json.NewEncoder(w).Encode(list)
	})
}

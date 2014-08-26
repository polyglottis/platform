// Package api contains the definition of all web services to extract data from the Polyglottis Platform.
package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/polyglottis/platform/content"
	"github.com/polyglottis/platform/language"
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

type resultList struct {
	ExtractCount int
	Results      []*result
}

type result struct {
	ExtractId content.ExtractId
	Slug      string
	Type      content.ExtractType
	Summaries []*summary
	Languages []language.Code
}

type summary struct {
	Title      string
	Summary    string
	Language   language.Code
	FlavorType content.FlavorType
	FlavorId   content.FlavorId
}

func (s *Server) ExtractSearch(r *http.Request) ([]byte, error) {
	return call(func(w io.Writer) error {
		query := r.URL.Query()
		q := &content.Query{
			LanguageA:   language.Code(query.Get("langA")),
			LanguageB:   language.Code(query.Get("langB")),
			ExtractType: content.ExtractType(query.Get("type")),
		}
		list, err := s.Content.ExtractsMatching(q)
		if err != nil {
			return err
		}

		results := make([]*result, 0)
		for _, id := range list {
			e, err := s.Content.GetExtract(id)
			if err != nil {
				if err != content.ErrNotFound {
					log.Println("Weird, we got this extract a second ago but it doesn't seem to exist:", id)
					continue
				} else {
					return err
				}
			}

			results = append(results, newResult(e, q))
		}
		return json.NewEncoder(w).Encode(newResultList(results))
	})
}

func newResultList(results []*result) *resultList {
	return &resultList{
		ExtractCount: len(results),
		Results:      results,
	}
}

func newResult(e *content.Extract, q *content.Query) *result {
	r := &result{
		ExtractId: e.Id,
		Slug:      e.UrlSlug,
		Type:      e.Type,
	}
	for lang, fByType := range e.Flavors {
		r.Languages = append(r.Languages, lang)
		if flavors, ok := fByType[content.Text]; ok {
			for _, f := range flavors {
				r.Summaries = append(r.Summaries, newSummary(lang, content.Text, f))
			}
		}
	}
	return r
}

func newSummary(lang language.Code, fType content.FlavorType, f *content.Flavor) *summary {
	s := &summary{
		Language:   lang,
		Summary:    f.Summary,
		FlavorType: fType,
		FlavorId:   f.Id,
	}
	if title := f.GetTitle(); title != nil {
		s.Title = title.Content
	}
	return s
}

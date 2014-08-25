// Package handle contains frontend handler functions.
package handle

import (
	"fmt"
	"net/http"

	"github.com/polyglottis/platform/backend"
	"github.com/polyglottis/platform/content"
	"github.com/polyglottis/platform/frontend"
)

type Worker struct {
	*backend.Engine
	Server          frontend.Server
	languageListSet bool
}

func NewWorker(engine *backend.Engine, server frontend.Server) *Worker {
	w := &Worker{
		Engine: engine,
		Server: server,
	}
	go w.fetchLanguageListPeriodically()
	return w
}

type Redirect struct {
	UrlStr string
	Code   int
}

func redirectToOther(urlStr string) *Redirect {
	return &Redirect{UrlStr: urlStr, Code: http.StatusSeeOther}
}

// Error is a hack: redirect is an error, thanks to this method
func (r *Redirect) Error() string {
	return fmt.Sprintf("Redirect [%d] to %s", r.Code, r.UrlStr)
}

func (w *Worker) readExtract(context *frontend.Context) (extract *content.Extract, err error) {
	slug := context.Vars["slug"]
	if len(slug) == 0 {
		err = content.ErrInvalidInput
		return
	}

	var id content.ExtractId
	id, err = w.Content.GetExtractId(slug)
	if err != nil {
		return
	}

	extract, err = w.Content.GetExtract(id)
	if err != nil {
		return
	}
	return
}

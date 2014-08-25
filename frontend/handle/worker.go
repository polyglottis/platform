package handle

import (
	"fmt"
	"net/http"

	"github.com/polyglottis/platform/backend"
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

package frontend

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/polyglottis/platform/backend"
	"github.com/polyglottis/platform/content"
	"github.com/polyglottis/platform/language"
	"github.com/polyglottis/platform/server"
)

type Worker struct {
	*backend.Engine
	Server Server
}

func NewWorker(engine *backend.Engine, frontendServer Server) *Worker {
	return &Worker{
		Engine: engine,
		Server: frontendServer,
	}
}

func (w *Worker) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/extract/edit/text", w.contextHandler(w.EditText))

	r.HandleFunc("/extract/{slug}", w.contextHandler(w.Extract))
	r.HandleFunc("/extract/{slug}/{language}", w.contextHandler(w.Flavor))

	r.HandleFunc("/", w.contextHandler(w.Server.Home))
	r.NotFoundHandler = http.HandlerFunc(w.contextHandlerCode(http.StatusNotFound, w.Server.NotFound))
}

func (w *Worker) contextHandler(f func(*Context) ([]byte, error)) func(http.ResponseWriter, *http.Request) {
	return w.contextHandlerCode(http.StatusOK, f)
}

func (worker *Worker) contextHandlerCode(code int, f func(*Context) ([]byte, error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				server.Recovered(r, w, rec)
			}
		}()
		context, err := ReadContext(r)
		if err == nil {
			bytes, err := f(context)
			switch {
			case err == nil:
				w.WriteHeader(code)
				w.Write(bytes)
			case err == content.ErrNotFound:
				worker.contextHandlerCode(http.StatusNotFound, worker.Server.NotFound)(w, r)
			default:
				server.InternalError(r, w, err)
			}
		} else {
			server.InternalError(r, w, err)
		}
	}
}

func (w *Worker) Extract(context *Context) ([]byte, error) {
	slug := context.Vars["slug"]
	if len(slug) == 0 {
		return nil, content.ErrNotFound
	}

	id, err := w.Content.GetExtractId(slug)
	if err != nil {
		return nil, err
	}

	return w.extract(context, id)
}

func (w *Worker) extract(context *Context, id content.ExtractId) ([]byte, error) {
	extract, err := w.Content.GetExtract(id)
	if err != nil {
		return nil, err
	}
	context.ExtractId = id
	return w.Server.Extract(context, extract)
}

func (w *Worker) Flavor(context *Context) ([]byte, error) {
	slug := context.Vars["slug"]
	lang := context.Vars["language"]
	if len(slug) == 0 {
		return nil, content.ErrNotFound
	}

	id, err := w.Content.GetExtractId(slug)
	if err != nil {
		return nil, err
	}

	langCode, err := w.Language.GetCode(lang)
	if err != nil {
		// language not found, fall back to extract
		return w.extract(context, id)
	}
	context.LanguageA = langCode

	extract, err := w.Content.GetExtract(id)
	if err != nil {
		return nil, err
	}
	context.ExtractId = id

	if fByType, ok := extract.Flavors[langCode]; ok {
		a := &FlavorTriple{}
		if audio, ok := fByType[content.Audio]; ok {
			a.Audio = audio[0]
		}
		if text, ok := fByType[content.Text]; ok {
			a.Text = text[0]
		}
		if transcript, ok := fByType[content.Transcript]; ok {
			a.Transcript = transcript[0]
		}
		return w.Server.Flavor(context, extract, a, &FlavorTriple{})
	}
	// flavor not found, fall back to extract
	return w.Server.Extract(context, extract)
}

func (w *Worker) EditText(context *Context) ([]byte, error) {
	id := content.ExtractId(context.Query.Get("id"))
	langA := context.Query.Get("a")
	langB := context.Query.Get("b")
	if len(id) == 0 || len(langA) == 0 {
		return nil, content.ErrInvalidInput
	}

	extract, err := w.Content.GetExtract(id)
	if err != nil {
		return nil, err
	}
	context.ExtractId = id

	langCodeA, err := w.Language.GetCode(langA)
	if err != nil {
		return nil, err
	}
	context.LanguageA = langCodeA

	var langCodeB language.Code
	if len(langB) != 0 {
		langCodeB, err = w.Language.GetCode(langB)
		if err != nil {
			return nil, err
		}
		context.LanguageB = langCodeB
	}

	if fByTypeA, ok := extract.Flavors[langCodeA]; ok {
		if textA, ok := fByTypeA[content.Text]; ok {
			var textB *content.Flavor
			if len(langCodeB) != 0 {
				if fByTypeB, ok := extract.Flavors[langCodeB]; ok {
					if tB, ok := fByTypeB[content.Text]; ok {
						textB = tB[0]
					}
				}
			}
			return w.Server.EditText(context, extract, textA[0], textB)
		}
	}
	return nil, content.ErrInvalidInput
}

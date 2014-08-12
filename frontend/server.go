package frontend

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/polyglottis/platform/backend"
	"github.com/polyglottis/platform/content"
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
	r.HandleFunc("/extract/{slug}", w.contextHandler(w.Extract))
	r.HandleFunc("/extract/{slug}/{flavor}", w.contextHandler(w.Flavor))

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
	return w.Server.Extract(context, extract)
}

func (w *Worker) Flavor(context *Context) ([]byte, error) {
	slug := context.Vars["slug"]
	flavor := context.Vars["flavor"]
	if len(slug) == 0 {
		return nil, content.ErrNotFound
	}

	id, err := w.Content.GetExtractId(slug)
	if err != nil {
		return nil, err
	}

	flavorIdInt, err := strconv.Atoi(flavor)
	if err != nil {
		// error parsing flavor, fall back to extract
		return w.extract(context, id)
	}
	flavorId := content.FlavorId(flavorIdInt)

	extract, err := w.Content.GetExtract(id)
	if err != nil {
		return nil, err
	}

	for _, f := range extract.Flavors {
		if f.Id == flavorId {
			return w.Server.Flavor(context, extract, f)
		}
	}
	// flavor not found, fall back to extract
	return w.Server.Extract(context, extract)
}

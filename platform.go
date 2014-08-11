// Package platform is the entry point to launche the Polyglottis Platform.
package platform

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/polyglottis/platform/api"
	"github.com/polyglottis/platform/backend"
	"github.com/polyglottis/platform/frontend"
)

type Configuration struct {
	Backend  *backend.Configuration
	Frontend frontend.Server
}

// Launch starts listening to addr and serving contents based on the configuration.
func Launch(addr string, c *Configuration) error {
	engine := backend.NewEngine(c.Backend)
	server := NewServer(addr)
	api.NewServer(engine).RegisterServices(server.Subrouter("/api"))
	frontend.NewWorker(engine, c.Frontend).RegisterRoutes(server.Router)
	return server.ListenAndServe()
}

type MainServer struct {
	Router *mux.Router
	http   *http.Server
}

func NewServer(addr string) *MainServer {
	r := mux.NewRouter()
	return &MainServer{
		Router: r,
		http: &http.Server{
			Addr:    addr,
			Handler: r,
		},
	}
}

func (s *MainServer) Subrouter(pathPrefix string) *mux.Router {
	return s.Router.PathPrefix(pathPrefix).Subrouter()
}

func (s *MainServer) ListenAndServe() error {
	return s.http.ListenAndServe()
}

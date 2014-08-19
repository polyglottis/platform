package frontend

import (
	"fmt"
	"net/http"

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
	r.HandleFunc("/user/signup", w.contextHandler(w.Server.SignUp)).Methods("GET")
	r.HandleFunc("/user/signin", w.contextHandler(w.Server.SignIn)).Methods("GET")
	r.HandleFunc("/user/signup", w.contextHandlerForm(w.SignUp)).Methods("POST")
	r.HandleFunc("/user/signin", w.contextHandlerForm(w.SignIn)).Methods("POST")
	r.HandleFunc("/user/signout", w.contextHandlerForm(w.SignOut)).Methods("GET")

	r.HandleFunc("/user/forgot_password", w.contextHandler(w.Server.ForgotPassword)).Methods("GET")
	r.HandleFunc("/user/forgot_password", w.contextHandlerForm(w.ForgotPassword)).Methods("POST")
	r.HandleFunc("/user/reset_password/{user}/{token}", w.contextHandler(w.GetResetPassword)).Methods("GET")
	r.HandleFunc("/user/reset_password/{user}/{token}", w.contextHandlerForm(w.ResetPassword)).Methods("POST")

	r.HandleFunc("/extract/edit/text", w.contextHandler(w.EditText))

	r.HandleFunc("/extract/{slug}", w.contextHandler(w.Extract))
	r.HandleFunc("/extract/{slug}/{language}", w.contextHandler(w.Flavor))

	r.HandleFunc("/", w.contextHandler(w.Server.Home))
	r.NotFoundHandler = http.HandlerFunc(w.contextHandlerCode(http.StatusNotFound, w.Server.NotFound))
}

func (w *Worker) contextHandler(f func(*Context) ([]byte, error)) func(http.ResponseWriter, *http.Request) {
	return w.contextHandlerCode(http.StatusOK, f)
}

func (w *Worker) contextHandlerForm(f func(*Context, *Session) ([]byte, error)) func(http.ResponseWriter, *http.Request) {
	return w.contextHandlerFull(http.StatusOK, f, true)
}

func (w *Worker) contextHandlerCode(code int, f func(*Context) ([]byte, error)) func(http.ResponseWriter, *http.Request) {
	return w.contextHandlerFull(code, func(c *Context, s *Session) ([]byte, error) {
		return f(c)
	}, false)
}

func (worker *Worker) contextHandlerFull(code int, f func(*Context, *Session) ([]byte, error), hasForm bool) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				server.Recovered(r, w, rec)
			}
		}()
		var err error
		var context *Context
		session := readSession(r, w)
		if hasForm {
			context, err = ReadContextWithForm(r, session)
		} else {
			context, err = ReadContext(r, session)
		}
		if err == nil {
			bytes, err := f(context, session)
			switch {
			case err == nil:
				w.WriteHeader(code)
				w.Write(bytes)
			case err == content.ErrNotFound:
				worker.contextHandlerCode(http.StatusNotFound, worker.Server.NotFound)(w, r)
			default:
				// hack to redirect: it is not an error to redirect, but it is handy to return a redirection in the error...
				if redir, ok := err.(*redirect); ok {
					http.Redirect(w, r, redir.urlStr, redir.code)
				} else {
					server.InternalError(r, w, err)
				}
			}
		} else {
			server.InternalError(r, w, err)
		}
	}
}

type redirect struct {
	urlStr string
	code   int
}

func redirectTo(urlStr string, code int) *redirect {
	return &redirect{urlStr: urlStr, code: code}
}

// Error is a hack: redirect is an error, thanks to this method
func (r *redirect) Error() string {
	return fmt.Sprintf("Redirect [%d] to %s", r.code, r.urlStr)
}

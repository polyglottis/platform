// Package router contains the route definitions for the frontend server.
package router

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/polyglottis/platform/backend"
	"github.com/polyglottis/platform/content"
	"github.com/polyglottis/platform/frontend"
	"github.com/polyglottis/platform/frontend/handle"
	"github.com/polyglottis/platform/server"
)

type Router struct {
	*handle.Worker
}

func NewRouter(engine *backend.Engine, server frontend.Server) *Router {
	return &Router{
		Worker: handle.NewWorker(engine, server),
	}
}

func (w *Router) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/user/signup", w.contextHandler(w.Server.SignUp)).Methods("GET")
	r.HandleFunc("/user/signin", w.contextHandler(w.Server.SignIn)).Methods("GET")
	r.HandleFunc("/user/signup", w.contextHandlerForm(w.SignUp)).Methods("POST")
	r.HandleFunc("/user/signin", w.contextHandlerForm(w.SignIn)).Methods("POST")
	r.HandleFunc("/user/signout", w.contextHandlerForm(w.SignOut)).Methods("GET")

	r.HandleFunc("/user/forgot_password", w.contextHandler(w.Server.ForgotPassword)).Methods("GET")
	r.HandleFunc("/user/forgot_password", w.contextHandlerForm(w.ForgotPassword)).Methods("POST")
	r.HandleFunc("/user/reset_password/{user}/{token}", w.contextHandler(w.GetResetPassword)).Methods("GET")
	r.HandleFunc("/user/reset_password/{user}/{token}", w.contextHandlerForm(w.ResetPassword)).Methods("POST")

	r.HandleFunc("/extract/edit/new", w.contextHandler(w.Server.NewExtract)).Methods("GET")
	r.HandleFunc("/extract/edit/new", w.contextHandlerForm(w.NewExtract)).Methods("POST")
	r.HandleFunc("/extract/edit/text", w.contextHandler(w.EditText))

	r.HandleFunc("/extract/{slug}", w.contextHandler(w.Extract))
	r.HandleFunc("/extract/{slug}/{language}", w.contextHandler(w.Flavor))

	r.HandleFunc("/", w.contextHandler(w.Server.Home))
	r.NotFoundHandler = http.HandlerFunc(w.contextHandlerCode(http.StatusNotFound, w.Server.NotFound))
}

func (w *Router) contextHandler(f func(*frontend.Context) ([]byte, error)) func(http.ResponseWriter, *http.Request) {
	return w.contextHandlerCode(http.StatusOK, f)
}

func (w *Router) contextHandlerForm(f func(*frontend.Context, *handle.Session) ([]byte, error)) func(http.ResponseWriter, *http.Request) {
	return w.contextHandlerFull(http.StatusOK, f, true)
}

func (w *Router) contextHandlerCode(code int, f func(*frontend.Context) ([]byte, error)) func(http.ResponseWriter, *http.Request) {
	return w.contextHandlerFull(code, forgetSession(f), false)
}

func (worker *Router) contextHandlerFull(code int, f func(*frontend.Context, *handle.Session) ([]byte, error), hasForm bool) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				server.Recovered(r, w, rec)
			}
		}()
		var err error
		var context *frontend.Context
		session := readSession(r, w)
		if hasForm {
			context, err = ReadContextWithForm(r, session)
		} else {
			context, err = ReadContext(r, session)
		}
		if err == nil {
			worker.do(&job{Code: code, HasForm: hasForm, Context: context, Session: session,
				F: f, W: w, R: r})
		} else {
			server.InternalError(r, w, err)
		}
	}
}

func forgetSession(f func(*frontend.Context) ([]byte, error)) func(*frontend.Context, *handle.Session) ([]byte, error) {
	return func(c *frontend.Context, s *handle.Session) ([]byte, error) {
		return f(c)
	}
}

type job struct {
	Code      int
	F         func(*frontend.Context, *handle.Session) ([]byte, error)
	HasForm   bool
	Session   *handle.Session
	Context   *frontend.Context
	W         http.ResponseWriter
	R         *http.Request
	secondTry bool
}

func (w *Router) do(job *job) {
	bytes, err := job.F(job.Context, job.Session)
	switch {
	case err == nil:
		job.W.WriteHeader(job.Code)
		job.W.Write(bytes)
	case err == content.ErrNotFound && !job.secondTry:
		job.Code = http.StatusNotFound
		job.F = forgetSession(w.Server.NotFound)
		job.secondTry = true
		w.do(job)
	default:
		// hack to redirect: it is not an error to redirect, but it is handy to return a redirection in the error...
		if redir, ok := err.(*handle.Redirect); ok {
			http.Redirect(job.W, job.R, redir.UrlStr, redir.Code)
		} else {
			if job.secondTry {
				server.InternalError(job.R, job.W, err)
			} else {
				log.Println("Frontend server error:", err)
				job.Code = http.StatusInternalServerError
				job.F = forgetSession(w.Server.Error)
				job.secondTry = true
				w.do(job)
			}
		}
	}
}

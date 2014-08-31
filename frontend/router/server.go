// Package router contains the route definitions for the frontend server.
package router

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/polyglottis/platform/backend"
	"github.com/polyglottis/platform/config"
	"github.com/polyglottis/platform/content"
	"github.com/polyglottis/platform/frontend"
	"github.com/polyglottis/platform/frontend/handle"
	"github.com/polyglottis/platform/server"
	"github.com/polyglottis/sitemap"
)

type Router struct {
	*handle.Worker
}

func NewRouter(engine *backend.Engine, server frontend.Server) *Router {
	return &Router{
		Worker: handle.NewWorker(engine, server),
	}
}

func (w *Router) RegisterRoutes(router *mux.Router) {
	c := config.Get()

	r := sitemap.NewRouter(router, c.Host, c.SitemapPath)

	r.Register("/user/signup").HandlerFunc(w.contextHandler(w.Server.SignUp)).Methods("GET")
	r.Register("/user/signin").HandlerFunc(w.contextHandler(w.Server.SignIn)).Methods("GET")
	r.HandleFunc("/user/signup", w.contextHandlerForm(w.SignUp)).Methods("POST")
	r.HandleFunc("/user/signin", w.contextHandlerForm(w.SignIn)).Methods("POST")
	r.HandleFunc("/user/signout", w.contextHandlerForm(w.SignOut)).Methods("GET")

	r.Register("/user/forgot_password").HandlerFunc(w.contextHandler(w.Server.ForgotPassword)).Methods("GET")
	r.HandleFunc("/user/forgot_password", w.contextHandlerForm(w.ForgotPassword)).Methods("POST")
	r.HandleFunc("/user/reset_password/{user}/{token}", w.contextHandlerSession(w.GetResetPassword)).Methods("GET")
	r.HandleFunc("/user/reset_password/{user}/{token}", w.contextHandlerForm(w.ResetPassword)).Methods("POST")

	r.Register("/extract/edit/new").HandlerFunc(w.contextHandler(w.Server.NewExtract)).Methods("GET")
	r.HandleFunc("/extract/edit/new", w.contextHandlerForm(w.NewExtract)).Methods("POST")
	r.HandleFunc("/extract/edit/new_flavor/{slug}",
		w.contextHandler(w.NewFlavor)).Methods("GET")
	r.HandleFunc("/extract/edit/new_flavor/{slug}",
		w.contextHandlerForm(w.NewFlavorPOST)).Methods("POST")
	r.HandleFunc("/extract/edit/text/{slug}", w.contextHandler(w.EditText)).Methods("GET")
	r.HandleFunc("/extract/edit/text/{slug}", w.contextHandlerForm(w.EditTextPOST)).Methods("POST")

	r.HandleFunc("/extract/{slug}", w.contextHandler(w.Extract))
	r.RegisterParam("/extract/{slug}/{language}", w.enumerateTexts).HandlerFunc(w.contextHandler(w.Flavor))

	r.Register("/").HandlerFunc(w.contextHandler(w.Server.Home))

	r.HandleSitemaps()
	go w.generateSitemapsEvery(r, 24*time.Hour)

	r.NotFoundHandler = http.HandlerFunc(w.contextHandlerCode(http.StatusNotFound, w.Server.NotFound))
}

func (w *Router) contextHandler(f func(*frontend.Context) ([]byte, error)) func(http.ResponseWriter, *http.Request) {
	return w.contextHandlerCode(http.StatusOK, f)
}

func (w *Router) contextHandlerSession(f func(*frontend.Context, *handle.Session) ([]byte, error)) func(http.ResponseWriter, *http.Request) {
	return w.contextHandlerFull(http.StatusOK, f, false)
}

func (w *Router) contextHandlerForm(f func(*frontend.Context, *handle.Session) ([]byte, error)) func(http.ResponseWriter, *http.Request) {
	return w.contextHandlerFull(http.StatusOK, f, true)
}

func (w *Router) contextHandlerCode(code int, f func(*frontend.Context) ([]byte, error)) func(http.ResponseWriter, *http.Request) {
	return w.contextHandlerFull(code, forgetSession(f), false)
}

func (router *Router) contextHandlerFull(code int, f func(*frontend.Context, *handle.Session) ([]byte, error), hasForm bool) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				server.Recovered(r, w, rec)
			}
		}()
		var err error
		var context *frontend.Context
		session := handle.NewSession(r, w)
		if hasForm {
			context, err = ReadContextWithForm(r, session)
		} else {
			context, err = ReadContext(r, session)
		}
		if err == nil {
			router.do(&job{Code: code, HasForm: hasForm, Context: context, Session: session,
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
				log.Printf("Internal Server Error while serving %s [%s]: %v", job.Context.Url, job.R.Method, err)
				job.Code = http.StatusInternalServerError
				job.F = forgetSession(w.Server.Error)
				job.secondTry = true
				w.do(job)
			}
		}
	}
}

func (r *Router) enumerateTexts(cb func(pairs ...string) error) error {
	extractList, err := r.Content.ExtractList()
	if err != nil {
		return err
	}

	for _, extract := range extractList {
		// get extract flavors
		e, err := r.Content.GetExtract(extract.Id)
		if err != nil {
			return err
		}

		for langCode, fByType := range e.Flavors {
			if _, ok := fByType[content.Text]; ok {
				cb("slug", e.UrlSlug, "language", string(langCode))
			}
		}
	}
	return nil
}

func (r *Router) generateSitemapsEvery(router *sitemap.Router, interval time.Duration) {
	for _ = range time.Tick(interval) {
		log.Println("Re-creating sitemap")
		_, err := router.GenerateSitemaps()
		if err != nil {
			log.Println("Failed to generate sitemaps:", err)
		}
	}
}

package handle

import (
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/sessions"

	"github.com/polyglottis/platform/frontend"
	"github.com/polyglottis/platform/i18n"
	"github.com/polyglottis/platform/user"
)

type Session struct {
	session  *sessions.Session
	defaults *sessions.Session
	r        *http.Request
	w        http.ResponseWriter
}

func newSession(s *sessions.Session, defaults *sessions.Session, r *http.Request, w http.ResponseWriter) *Session {
	return &Session{
		session:  s,
		defaults: defaults,
		r:        r,
		w:        w,
	}
}

func (s *Session) SetAccount(a *user.Account) {
	s.session.Values["account"] = a
}

func (s *Session) GetAccount() *user.Account {
	if u, ok := s.session.Values["account"]; ok && u != nil {
		if account, ok := u.(*user.Account); ok {
			return account
		} else {
			log.Println("Unable to decode user account: did user.Account change recently?")
		}
	}
	return nil
}

func (s *Session) RemoveAccount() {
	delete(s.session.Values, "account")
}

func (s *Session) Save() error {
	return s.session.Save(s.r, s.w)
}

func (s *Session) SaveFlashErrors(errMap frontend.ErrorMap) {
	s.session.AddFlash(errMap)
	s.Save()
}

func (s *Session) SaveFlashError(msg i18n.Key) {
	s.SaveFlashErrors(frontend.ErrorMap{
		"FORM": msg,
	})
	s.Save()
}

func (s *Session) ReadFlashErrors() frontend.ErrorMap {
	if flashes := s.session.Flashes(); len(flashes) != 0 {
		defer s.Save()
		if errMap, ok := flashes[0].(frontend.ErrorMap); ok {
			return errMap
		} else {
			log.Println("Flash message is not an error map:", flashes[0])
		}
		if len(flashes) > 1 {
			log.Println("Session with multiple flash messages:", flashes)
		}
		log.Println("No flashes...")
	}
	return nil
}

func (s *Session) SaveDefaults(values url.Values) {
	s.defaults.Options = &sessions.Options{
		Path:   s.r.URL.String(),
		MaxAge: 3600,
	}
	s.defaults.Values["def"] = values
	s.defaults.Save(s.r, s.w)
}

func (s *Session) GetDefaults() url.Values {
	if def, ok := s.defaults.Values["def"]; ok {
		if defaults, ok := def.(url.Values); ok {
			return defaults
		}
	}
	return nil
}

func (s *Session) ClearDefaults() {
	s.defaults.Options = &sessions.Options{
		Path:   s.r.URL.String(),
		MaxAge: 0,
	}
	s.defaults.Save(s.r, s.w)
}

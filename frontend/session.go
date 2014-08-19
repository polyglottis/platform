package frontend

import (
	"encoding/gob"
	"log"
	"net/http"

	"github.com/gorilla/sessions"

	"github.com/polyglottis/platform/user"
)

func init() {
	gob.Register(&user.Account{})
}

var sessionStore = sessions.NewCookieStore([]byte("something-very-secret"))

type Session struct {
	*sessions.Session
	r *http.Request
	w http.ResponseWriter
}

func (s *Session) SetAccount(a *user.Account) {
	s.Values["account"] = a
}

func (s *Session) GetAccount() *user.Account {
	if u, ok := s.Values["account"]; ok && u != nil {
		if account, ok := u.(*user.Account); ok {
			return account
		} else {
			log.Println("Unable to decode user account: did user.Account change recently?")
		}
	}
	return nil
}

func (s *Session) RemoveAccount() {
	delete(s.Values, "account")
}

func (s *Session) Save() error {
	return s.Session.Save(s.r, s.w)
}

func readSession(r *http.Request, w http.ResponseWriter) *Session {
	s, err := sessionStore.Get(r, "user")
	if err != nil {
		log.Println("Unable to decode old session:", err)
	}
	return &Session{
		Session: s,
		r:       r,
		w:       w,
	}
}

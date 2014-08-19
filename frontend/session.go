package frontend

import (
	"encoding/gob"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"

	"github.com/polyglottis/platform/config"
	"github.com/polyglottis/platform/user"
)

func init() {
	gob.Register(&user.Account{})

	keyFile := config.Get().SessionKeyPath
	var keyPairs []*keyPair
	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		pairs := []*keyPair{newKeyPair()}
		toSave, err := json.Marshal(pairs)
		if err != nil {
			log.Fatalf("Unable to marshal session key pairs: %v", err)
		}

		f, err := os.Create(keyFile)
		if err != nil {
			log.Fatalf("Unable to create session key file at %s: %v", keyFile, err)
		}

		_, err = f.Write(toSave)
		if err != nil {
			log.Fatalf("Unable to store session keys in file %s: %v", keyFile, err)
		}
		log.Println("Session keys file created")
	} else {
		f, err := ioutil.ReadFile(keyFile)
		if err != nil {
			if abs, err := filepath.Abs(keyFile); err == nil {
				log.Fatalf("Unable to read session keys file [%s]: %v", abs, err)
			} else {
				log.Fatalf("Unable to read session keys file: %v", err)
			}
		}
		err = json.Unmarshal(f, &keyPairs)
		if err != nil {
			log.Fatalf("Error parsing session keys file: %v", err)
		}
	}

	keys := make([][]byte, 2*len(keyPairs))
	for i, kp := range keyPairs {
		keys[2*i] = kp.AuthKey
		keys[2*i+1] = kp.EncryptKey
	}
	sessionStore = sessions.NewCookieStore(keys...)
}

type keyPair struct {
	AuthKey    []byte
	EncryptKey []byte
	Time       time.Time // store time to allow key rotation (TODO)
}

func newKeyPair() *keyPair {
	return &keyPair{
		AuthKey:    securecookie.GenerateRandomKey(64),
		EncryptKey: securecookie.GenerateRandomKey(32),
		Time:       time.Now(),
	}
}

var sessionStore *sessions.CookieStore

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
		s, err = sessionStore.New(r, "user")
		if err != nil {
			log.Println("Unable to create new session: is there a problem with the session keys file?", err)
		}
	}
	return &Session{
		Session: s,
		r:       r,
		w:       w,
	}
}

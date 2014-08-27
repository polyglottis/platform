package handle

import (
	"encoding/gob"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"

	"github.com/polyglottis/platform/config"
	"github.com/polyglottis/platform/frontend"
	"github.com/polyglottis/platform/user"
)

func init() {
	gob.Register(&user.Account{})
	gob.Register(frontend.ErrorMap{})
	gob.Register(url.Values{})

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
		defer f.Close()

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
	sessionStore.Options.Path = "/"
	sessionStore.Options.MaxAge = 86400 * 7 // 1 week
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

func readSession(key string, r *http.Request) *sessions.Session {
	s, err := sessionStore.Get(r, key)
	if err != nil {
		log.Printf("Unable to decode old %s session: %v", key, err)
		s, err = sessionStore.New(r, key)
		if err != nil {
			log.Printf("Unable to create new %s session: is there a problem with the session keys file? %v", key, err)
		}
	}
	return s
}

func NewSession(r *http.Request, w http.ResponseWriter) *Session {
	s := readSession("user", r)
	def := readSession("defaults", r)
	return newSession(s, def, r, w)
}

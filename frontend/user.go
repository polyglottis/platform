package frontend

import (
	"crypto/rand"
	"math/big"
	mathRand "math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/schema"

	"github.com/polyglottis/platform/i18n"
	"github.com/polyglottis/platform/user"
	"github.com/polyglottis/platform/user/password"
)

type SignUpArgs struct {
	User            string
	Email           string
	Password        string
	PasswordConfirm string
}

func (a *SignUpArgs) CleanUp() {
	a.User = strings.TrimSpace(a.User)
	a.Email = strings.TrimSpace(a.Email)
}

var decoder = schema.NewDecoder()

func (w *Worker) SignUp(context *Context, session *Session) ([]byte, error) {
	args := new(SignUpArgs)
	err := decoder.Decode(args, context.Form)
	if err != nil {
		return nil, err
	}
	args.CleanUp()

	errors := make(map[string]i18n.Key)

	if valid, msg := user.ValidName(args.User); valid {
		_, err = w.User.GetAccount(user.Name(args.User))
		if err != nil && err != user.AccountNotFound {
			return nil, err
		}
		if err == nil {
			errors["User"] = i18n.Key("Username is already taken")
		}
	} else {
		errors["User"] = msg
	}

	if !validEmail(args.Email) {
		errors["Email"] = i18n.Key("Invalid email address")
	}

	if valid, msg := user.ValidPassword(args.Password); valid {
		if args.Password != args.PasswordConfirm {
			errors["Password"] = i18n.Key("Password doesn't match the confirmation")
		}
	} else {
		errors["Password"] = msg
	}

	if len(errors) != 0 {
		context.Errors = errors
		context.Defaults = url.Values{}
		context.Defaults.Set("User", args.User)
		context.Defaults.Set("Email", args.Email)
		return w.Server.SignUp(context)
	}

	hash, err := password.Hash(args.Password)
	if err != nil {
		return nil, err
	}

	a, err := w.User.NewAccount(&user.NewAccountRequest{
		Name:         user.Name(args.User),
		Email:        args.Email,
		MainLanguage: context.Locale,
		PasswordHash: hash,
	})
	if err != nil {
		return nil, err
	}

	session.SetAccount(a)
	err = session.Save()
	if err != nil {
		return nil, err
	}

	return nil, redirectTo("/", http.StatusSeeOther)
}

var emailRegex = regexp.MustCompile(`(?i)^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,4}$`)

func validEmail(email string) bool {
	return emailRegex.MatchString(email)
}

type SignInArgs struct {
	User     string
	Password string
}

func (a *SignInArgs) CleanUp() {
	a.User = strings.TrimSpace(a.User)
}

func (w *Worker) SignIn(context *Context, session *Session) ([]byte, error) {
	args := new(SignInArgs)
	err := decoder.Decode(args, context.Form)
	if err != nil {
		return nil, err
	}
	args.CleanUp()

	a, err := w.User.GetAccount(user.Name(args.User))
	if err != nil && err != user.AccountNotFound {
		return nil, err
	}
	if err == user.AccountNotFound || !password.Check(args.Password, a) {
		context.Errors = map[string]i18n.Key{
			"FORM": i18n.Key("Incorrect username or password."),
		}
		context.Defaults = url.Values{}
		context.Defaults.Set("User", args.User)
		sleep()
		return w.Server.SignIn(context)
	}

	session.SetAccount(a)
	err = session.Save()
	if err != nil {
		return nil, err
	}

	returnTo := context.Query.Get("return_to")
	if len(returnTo) == 0 {
		returnTo = "/"
	}
	return nil, redirectTo(returnTo, http.StatusSeeOther)
}

func sleep() {
	var i int64
	var limit int64 = 1000
	bigInt, err := rand.Int(rand.Reader, big.NewInt(limit))
	if err == nil {
		i = bigInt.Int64()
	} else {
		i = mathRand.Int63n(limit)
	}
	t := 1*time.Second + time.Duration(i)*time.Millisecond
	time.Sleep(t)
}

func (w *Worker) SignOut(context *Context, session *Session) ([]byte, error) {
	session.RemoveAccount()
	err := session.Save()
	if err != nil {
		return nil, err
	}

	return nil, redirectTo("/", http.StatusSeeOther)
}

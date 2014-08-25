package frontend

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/polyglottis/platform/i18n"
	"github.com/polyglottis/platform/user"
	"github.com/polyglottis/platform/user/password"
)

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

func (w *Worker) SignOut(context *Context, session *Session) ([]byte, error) {
	session.RemoveAccount()
	err := session.Save()
	if err != nil {
		return nil, err
	}

	return nil, redirectTo("/", http.StatusSeeOther)
}

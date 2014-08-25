package handle

import (
	"strings"

	"github.com/polyglottis/platform/frontend"
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

func (w *Worker) SignIn(context *frontend.Context, session *Session) ([]byte, error) {
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
		sleep()
		session.SaveFlashError("Incorrect username or password.")
		session.SaveDefault("User", args.User)
		return nil, redirectToOther(context.Url)
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
	session.ClearDefaults()
	return nil, redirectToOther(returnTo)
}

func (w *Worker) SignOut(context *frontend.Context, session *Session) ([]byte, error) {
	session.RemoveAccount()
	err := session.Save()
	if err != nil {
		return nil, err
	}

	return nil, redirectToOther("/")
}

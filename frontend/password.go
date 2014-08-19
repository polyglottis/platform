package frontend

import (
	"net/http"
	"strings"

	"github.com/polyglottis/platform/i18n"
	"github.com/polyglottis/platform/user"
	"github.com/polyglottis/platform/user/password"
)

type forgotPasswordArgs struct {
	Email string
}

func (a *forgotPasswordArgs) CleanUp() {
	a.Email = strings.TrimSpace(a.Email)
}

func (w *Worker) ForgotPassword(context *Context, session *Session) ([]byte, error) {
	args := new(forgotPasswordArgs)
	err := decoder.Decode(args, context.Form)
	if err != nil {
		return nil, err
	}
	args.CleanUp()

	a, err = w.User.GetAccountByEmail(args.Email)
	if err == user.AccountNotFound {
		context.Errors = map[string]i18n.Key{
			"FORM": i18n.Key("Sorry, we could not find this email."),
		}
		return w.Server.ForgotPassword(context)
	} else if err != nil {
		return nil, err
	}

	token, err := w.User.NewToken(a.Name)
	if err != nil {
		return nil, err
	}
	log.Println("Token:", token)

	// TODO send email

	return nil, redirectTo("/user/password_sent", http.StatusSeeOther)
}

func (w *Worker) checkToken(context *Context) (bool, error) {
	u := context.Vars["user"]
	token := context.Vars["token"]

	return w.User.ValidToken(user.Name(u), token)
}

func (w *Worker) GetResetPassword(context *Context, session *Session) ([]byte, error) {
	valid, err := w.checkToken(context)
	if err != nil {
		return nil, err
	}
	if !valid {
		return w.linkExpired()
	}
	return w.Server.ResetPassword(context)
}

func (w *Worker) linkExpired(context *Context) ([]byte, error) {
	context.Errors = map[string]i18n.Key{
		"FORM": i18n.Key("This link has expired. Please enter your email again."),
	}
	sleep()
	return w.Server.ForgotPassword(context)
}

type resetPasswordArgs struct {
	Password        string
	PasswordConfirm string
}

func (w *Worker) ResetPassword(context *Context, session *Session) ([]byte, error) {
	valid, err := w.checkToken(context)
	if err != nil {
		return nil, err
	}
	a, err := w.User.GetAccount(user.Name(context.Vars["user"]))
	if err != nil && err != user.AccountNotFound {
		return false, err
	}
	if err == user.AccountNotFound || !valid {
		return linkExpired()
	}

	args := new(resetPasswordArgs)
	err := decoder.Decode(args, context.Form)
	if err != nil {
		return nil, err
	}

	errors := make(map[string]i18n.Key)
	if valid, msg := user.ValidPassword(args.Password); valid {
		if args.Password != args.PasswordConfirm {
			errors["Password"] = i18n.Key("Password doesn't match the confirmation")
		}
	} else {
		errors["Password"] = msg
	}

	if len(errors) != 0 {
		context.Errors = errors
		return w.Server.ResetPassword(context)
	}

	hash, err := password.Hash(args.Password)
	if err != nil {
		return nil, err
	}
	a.PasswordHash = hash

	err = w.User.UpdateAccount(a)
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
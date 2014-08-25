package handle

import (
	"github.com/polyglottis/platform/frontend"
	"github.com/polyglottis/platform/i18n"
	"github.com/polyglottis/platform/user"
	"github.com/polyglottis/platform/user/password"
)

func (w *Worker) checkToken(context *frontend.Context) (bool, error) {
	u := context.Vars["user"]
	token := context.Vars["token"]

	return w.User.ValidToken(user.Name(u), token)
}

func (w *Worker) GetResetPassword(context *frontend.Context, session *Session) ([]byte, error) {
	valid, err := w.checkToken(context)
	if err != nil {
		return nil, err
	}
	if !valid {
		return w.linkExpired(session)
	}
	return w.Server.ResetPassword(context)
}

func (w *Worker) linkExpired(session *Session) ([]byte, error) {
	sleep()
	session.SaveFlashError("This link has expired. Please enter your email again.")
	return nil, redirectToOther("/user/forgot_password")
}

type resetPasswordArgs struct {
	Password        string
	PasswordConfirm string
}

func (w *Worker) ResetPassword(context *frontend.Context, session *Session) ([]byte, error) {
	valid, err := w.checkToken(context)
	if err != nil {
		return nil, err
	}
	a, err := w.User.GetAccount(user.Name(context.Vars["user"]))
	if err != nil && err != user.AccountNotFound {
		return nil, err
	}
	if err == user.AccountNotFound || !valid {
		return w.linkExpired(session)
	}

	args := new(resetPasswordArgs)
	err = decoder.Decode(args, context.Form)
	if err != nil {
		return nil, err
	}

	errors := make(frontend.ErrorMap)
	if valid, msg := user.ValidPassword(args.Password); valid {
		if args.Password != args.PasswordConfirm {
			errors["Password"] = i18n.Key("Password doesn't match the confirmation")
		}
	} else {
		errors["Password"] = msg
	}

	if len(errors) != 0 {
		session.SaveFlashErrors(errors)
		return nil, redirectToOther(context.Url)
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

	err = w.User.DeleteToken(user.Name(context.Vars["user"]), context.Vars["token"])
	if err != nil {
		return nil, err
	}

	session.SetAccount(a)
	err = session.Save()
	if err != nil {
		return nil, err
	}

	return nil, redirectToOther("/")
}

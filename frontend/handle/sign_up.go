package handle

import (
	"net/url"
	"strings"

	"github.com/polyglottis/platform/frontend"
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

func (w *Worker) SignUp(context *frontend.Context, session *Session) ([]byte, error) {
	args := new(SignUpArgs)
	err := decoder.Decode(args, context.Form)
	if err != nil {
		return nil, err
	}
	args.CleanUp()

	errors := make(frontend.ErrorMap)

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
		defaults := url.Values{}
		defaults.Set("User", args.User)
		defaults.Set("Email", args.Email)
		session.SaveDefaults(defaults)
		session.SaveFlashErrors(errors)
		return nil, redirectToOther(context.Url)
	}

	hash, err := password.Hash(args.Password)
	if err != nil {
		return nil, err
	}

	a, err := w.User.NewAccount(&user.NewAccountRequest{
		Name:         user.Name(args.User),
		Email:        args.Email,
		UILocale:     context.Locale,
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

	session.ClearDefaults()
	return nil, redirectToOther("/")
}

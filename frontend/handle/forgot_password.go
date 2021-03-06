package handle

import (
	"strings"

	"github.com/polyglottis/platform/config"
	"github.com/polyglottis/platform/frontend"
	"github.com/polyglottis/platform/user"
)

type forgotPasswordArgs struct {
	Email string
}

func (a *forgotPasswordArgs) CleanUp() {
	a.Email = strings.TrimSpace(a.Email)
}

func (w *Worker) ForgotPassword(context *frontend.Context, session *Session) ([]byte, error) {
	args := new(forgotPasswordArgs)
	err := decoder.Decode(args, context.Form)
	if err != nil {
		return nil, err
	}
	args.CleanUp()

	a, err := w.User.GetAccountByEmail(args.Email)
	if err == user.AccountNotFound {
		sleep()
		session.SaveFlashError("Sorry, we could not find this email.")
		return nil, redirectToOther("/user/forgot_password")
	} else if err != nil {
		return nil, err
	}

	token, err := w.User.NewToken(a.Name)
	if err != nil {
		return nil, err
	}

	message, err := w.Server.PasswordResetEmail(context, a, token)
	if err != nil {
		return nil, err
	}

	err = config.Get().MailUser.SendMail("support@polyglottis.org", message, args.Email)
	if err != nil {
		return nil, err
	}

	context.Email = args.Email
	return w.Server.PasswordSent(context)
}

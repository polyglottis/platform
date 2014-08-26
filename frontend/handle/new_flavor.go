package handle

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/polyglottis/platform/content"
	"github.com/polyglottis/platform/frontend"
	"github.com/polyglottis/platform/i18n"
)

type newFlavorArgs struct {
	Language     string
	MainComment  string
	OtherComment string
}

func (a *newFlavorArgs) CleanUp() {
	a.MainComment = strings.TrimSpace(a.MainComment)
	a.OtherComment = strings.TrimSpace(a.OtherComment)
}

func (w *Worker) NewFlavorPOST(context *frontend.Context, session *Session) ([]byte, error) {
	errors := make(frontend.ErrorMap)
	defaults := url.Values{}

	if !context.LoggedIn() {
		errors["FORM"] = i18n.Key("You must sign in to perform this action.")
	}

	e, err := w.readExtract(context)
	if err != nil {
		return nil, err
	}

	args := new(newFlavorArgs)
	err = decoder.Decode(args, context.Form)
	if err != nil {
		return nil, err
	}
	args.CleanUp()

	langCode, err := w.Language.GetCode(args.Language)
	if err == nil {
		defaults.Set("Language", args.Language)
	} else {
		defaults.Set("Language", "")
		errors["Language"] = i18n.Key("Please select an option.")
	}

	var mainComment string
	var otherFlavor *content.Flavor
	if fByType, ok := e.Flavors[langCode]; ok {
		if flavors, ok := fByType[content.Text]; ok {
			if len(flavors) > 0 {
				// MainComment only necessary if there are already other flavors.
				if valid, msg := content.ValidLanguageComment(args.MainComment); valid {
					mainComment = args.MainComment
				} else {
					errors["MainComment"] = msg
				}
			}
			if len(flavors) == 1 {
				// OtherComment only necessary if there is exactly one other flavor.
				if valid, msg := content.ValidLanguageComment(args.OtherComment); valid {
					otherFlavor = flavors[0]
					otherFlavor.LanguageComment = args.OtherComment
				} else {
					errors["OtherComment"] = msg
				}
			}
		}
	}

	if len(errors) != 0 {
		session.SaveFlashErrors(errors)
		defaults.Set("MainComment", args.MainComment)
		defaults.Set("OtherComment", args.OtherComment)
		session.SaveDefaults(defaults)
		return nil, redirectToOther(context.Url)
	}

	if otherFlavor != nil {
		err = w.Content.UpdateFlavor(context.User, otherFlavor)
		if err != nil {
			return nil, err
		}
	}

	f := &content.Flavor{
		ExtractId:       e.Id,
		Type:            content.Text,
		Language:        langCode,
		LanguageComment: mainComment,
	}

	err = w.Content.NewFlavor(context.User, f)
	if err != nil {
		return nil, err
	}

	session.ClearDefaults()
	context.Query.Set("a", args.Language)
	context.Query.Set("at", strconv.Itoa(int(f.Id)))
	return nil, redirectToOther("/extract/edit/text/" + e.UrlSlug + "?" + context.Query.Encode())
}

func (w *Worker) NewFlavor(context *frontend.Context) ([]byte, error) {
	extract, err := w.readExtract(context)
	if err != nil {
		return nil, err
	}
	// avoid sending too much over rpc
	// TODO do not read the blocks at all in the database
	for _, fByType := range extract.Flavors {
		for _, flavors := range fByType {
			for _, f := range flavors {
				f.Blocks = nil
			}
		}
	}

	return w.Server.NewFlavor(context, extract)
}

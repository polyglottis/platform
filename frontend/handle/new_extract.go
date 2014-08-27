package handle

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/polyglottis/platform/content"
	"github.com/polyglottis/platform/frontend"
	"github.com/polyglottis/platform/i18n"
)

type newExtractArgs struct {
	Slug        string
	ExtractType string
	Language    string
	Title       string
	Summary     string
	Text        string
}

func (a *newExtractArgs) CleanUp() {
	a.Slug = content.NormalizeSlug(a.Slug)
	a.Title = strings.TrimSpace(a.Title)
	a.Summary = strings.TrimSpace(a.Summary)
	a.Text = strings.TrimSpace(a.Text)
}

func (w *Worker) NewExtract(context *frontend.Context, session *Session) ([]byte, error) {
	errors := make(frontend.ErrorMap)
	defaults := url.Values{}

	if !context.LoggedIn() {
		errors["FORM"] = i18n.Key("You must sign in to perform this action.")
	}

	args := new(newExtractArgs)
	err := decoder.Decode(args, context.Form)
	if err != nil {
		return nil, err
	}
	args.CleanUp()

	defaults.Set("ExtractType", args.ExtractType)
	defaults.Set("Language", args.Language)

	if valid, msg := content.ValidSlug(args.Slug); valid {
		_, err = w.Content.GetExtractId(args.Slug)
		if err == nil {
			errors["Slug"] = i18n.Key("This url slug is already taken")
		} else if err != content.ErrNotFound {
			return nil, err
		}
	} else {
		errors["Slug"] = msg
	}

	if !content.ValidExtractType(content.ExtractType(args.ExtractType)) {
		defaults.Set("ExtractType", "")
		errors["ExtractType"] = i18n.Key("Please select an option.")
	}

	langCode, err := w.Language.GetCode(args.Language)
	if err != nil {
		defaults.Set("Language", "")
		errors["Language"] = i18n.Key("Please select an option.")
	}

	if valid, msg := content.ValidTitle(args.Title); !valid {
		errors["Title"] = msg
	}

	if valid, msg := content.ValidSummary(args.Summary); !valid {
		errors["Summary"] = msg
	}

	if len(args.Text) == 0 {
		errors["Text"] = i18n.Key("Please enter your extract.")
	}

	if len(errors) != 0 {
		session.SaveFlashErrors(errors)
		defaults.Set("Slug", args.Slug)
		defaults.Set("Title", args.Title)
		defaults.Set("Summary", args.Summary)
		defaults.Set("Text", args.Text)
		session.SaveDefaults(defaults)
		return nil, redirectToOther(context.Url)
	}

	e := &content.Extract{
		UrlSlug: args.Slug,
		Type:    content.ExtractType(args.ExtractType),
		Flavors: content.FlavorMap{
			langCode: content.FlavorByType{
				content.Text: []*content.Flavor{{
					Summary: args.Summary,
					Blocks:  getBlocks(args.Title, args.Text),
				}},
			},
		},
	}
	err = w.Content.NewExtract(context.User, e)
	if err != nil {
		return nil, err
	}

	session.ClearDefaults()
	return nil, redirectToOther("/extract/" + args.Slug + "/" + string(langCode))
}

var blockRegexp = regexp.MustCompile(`\s*\n\s*\n\s*`)
var unitRegexp = regexp.MustCompile(`\s*\n\s*`)

func getBlocks(title, text string) content.BlockSlice {
	c := content.BlockSlice{content.UnitSlice{{
		ContentType: content.TypeText,
		Content:     title,
	}}}

	textBlocks := blockRegexp.Split(text, -1)
	for _, b := range textBlocks {
		units := unitRegexp.Split(b, -1)
		block := make(content.UnitSlice, len(units))
		for i, u := range units {
			block[i] = &content.Unit{
				ContentType: content.TypeText,
				Content:     u,
			}
		}
		c = append(c, block)
	}
	return c
}

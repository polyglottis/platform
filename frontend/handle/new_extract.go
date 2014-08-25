package handle

import (
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"

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
	a.Slug = strings.TrimSpace(a.Slug)
	a.Title = strings.TrimSpace(a.Title)
	a.Summary = strings.TrimSpace(a.Summary)
	a.Text = strings.TrimSpace(a.Text)
}

func (w *Worker) NewExtract(context *frontend.Context, session *Session) ([]byte, error) {
	errors := make(map[string]i18n.Key)
	context.Defaults = url.Values{}

	if len(context.User) == 0 {
		errors["FORM"] = i18n.Key("You must be logged in to perform this action.")
	}

	args := new(newExtractArgs)
	err := decoder.Decode(args, context.Form)
	if err != nil {
		return nil, err
	}
	args.CleanUp()

	if valid, msg := content.ValidSlug(args.Slug); !valid {
		errors["Slug"] = msg
	}

	if !content.ValidExtractType(content.ExtractType(args.ExtractType)) {
		context.Defaults.Set("ExtractType", "")
		errors["ExtractType"] = i18n.Key("Please select one option.")
	}

	langCode, err := w.Language.GetCode(args.Language)
	if err != nil {
		context.Defaults.Set("Language", "")
		errors["Language"] = i18n.Key("Please select one option.")
	}

	if len(args.Title) == 0 {
		errors["Title"] = i18n.Key("Please enter a title.")
	}

	summaryLength := utf8.RuneCountInString(args.Summary)
	if summaryLength < 10 {
		errors["Summary"] = i18n.Key("Please enter a longer summary.")
	}
	if summaryLength > 150 {
		errors["Summary"] = i18n.Key("This summary is too long (maximum 150 characters).")
	}

	if len(args.Text) == 0 {
		errors["Text"] = i18n.Key("Please enter your extract.")
	}

	if len(errors) != 0 {
		context.Errors = errors
		context.Defaults.Set("Slug", args.Slug)
		context.Defaults.Set("Title", args.Title)
		context.Defaults.Set("Summary", args.Summary)
		context.Defaults.Set("Text", args.Text)
		return w.Server.NewExtract(context)
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
	return nil, redirectToOther("/extract/" + args.Slug)
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

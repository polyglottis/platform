package handle

import (
	"log"
	"net/url"
	"strings"

	"github.com/polyglottis/platform/content"
	"github.com/polyglottis/platform/frontend"
	"github.com/polyglottis/platform/i18n"
	"github.com/polyglottis/platform/language"
)

func (w *Worker) EditText(context *frontend.Context) ([]byte, error) {
	extract, langCodeA, langCodeB, err := w.readExtractAndLanguages(context)
	if err != nil {
		return nil, err
	}

	if fByTypeA, ok := extract.Flavors[langCodeA]; ok {
		a := newFlavorTriple(fByTypeA, context, "a")
		b := &frontend.FlavorTriple{}
		if len(langCodeB) != 0 {
			if fByTypeB, ok := extract.Flavors[langCodeB]; ok {
				b = newFlavorTriple(fByTypeB, context, "b")
			}
		}
		return w.Server.EditText(context, extract, a.Text, b.Text)
	}
	return nil, content.ErrInvalidInput
}

func (w *Worker) readExtractAndLanguages(context *frontend.Context) (extract *content.Extract, langCodeA, langCodeB language.Code, err error) {
	extract, err = w.readExtract(context)
	if err != nil {
		return
	}

	langA := context.Query.Get("a")
	langB := context.Query.Get("b")

	langCodeA, err = w.Language.GetCode(langA)
	if err != nil {
		return
	}

	if len(langB) != 0 {
		langCodeB, err = w.Language.GetCode(langB)
		if err != nil {
			return
		}
	}
	return
}

type editTextArgs struct {
	Summary string
	Title   string
	Blocks  []*block
}
type block struct {
	Units []*unit
}
type unit struct {
	BlockId int
	UnitId  int
	Content string
}

func (a *editTextArgs) CleanUp() {
	a.Summary = strings.TrimSpace(a.Summary)
	a.Title = strings.TrimSpace(a.Title)
	for _, block := range a.Blocks {
		for _, u := range block.Units {
			u.Content = strings.TrimSpace(u.Content)
		}
	}
}

func (a *editTextArgs) ContentBlocks() content.BlockSlice {
	blocks := make(content.BlockSlice, 0, len(a.Blocks))
	for _, b := range a.Blocks {
		block := make(content.UnitSlice, 0, len(b.Units))
		for _, u := range b.Units {
			block = append(block, &content.Unit{
				ContentType: content.TypeText,
				Content:     u.Content,
				BlockId:     content.BlockId(u.BlockId),
				Id:          content.UnitId(u.UnitId),
			})
		}
		if len(block) != 0 {
			blocks = append(blocks, block)
		}
	}
	return blocks
}

func (w *Worker) EditTextPOST(context *frontend.Context, session *Session) ([]byte, error) {
	errors := make(frontend.ErrorMap)

	if !context.LoggedIn() {
		errors["FORM"] = i18n.Key("You must sign in to perform this action.")
	}

	e, langCodeA, langCodeB, err := w.readExtractAndLanguages(context)
	if err != nil {
		return nil, err
	}

	var langCode language.Code
	var flavorId content.FlavorId
	if context.IsFocusOnA() {
		langCode = langCodeA
		flavorId, err = context.FlavorId("at")
	} else {
		langCode = langCodeB
		flavorId, err = context.FlavorId("bt")
	}
	if err != nil {
		return nil, err
	}

	args := new(editTextArgs)
	err = decoder.Decode(args, context.Form)
	if err != nil {
		return nil, err
	}
	args.CleanUp()

	if valid, msg := content.ValidSummary(args.Summary); !valid {
		errors["Summary"] = msg
	}

	if valid, msg := content.ValidTitle(args.Title); !valid {
		errors["Title"] = msg
	}

	log.Println("Lines", args.Blocks)

	if len(errors) != 0 {
		session.SaveFlashErrors(errors)
		defaults := url.Values{}
		defaults.Set("Summary", args.Summary)
		defaults.Set("Title", args.Title)
		// TODO set defaults for each (modified) line!
		session.SaveDefaults(defaults)
		return nil, redirectToOther(context.Url)
	}

	f, err := e.GetFlavor(langCode, content.Text, flavorId)
	if err != nil {
		return nil, err
	}

	if f.Summary != args.Summary {
		f.Summary = args.Summary
		err = w.Content.UpdateFlavor(context.User, f)
		if err != nil {
			return nil, err
		}
	}

	// compare old and new, record updates
	toUpdate := make([]*content.Unit, 0)
	curTitle := f.GetTitle()
	if curTitle == nil {
		toUpdate = append(toUpdate, &content.Unit{
			BlockId:     1,
			Id:          1,
			ContentType: content.TypeText,
			Content:     args.Title,
		})
	} else if curTitle.Content != args.Title {
		curTitle.Content = args.Title
		toUpdate = append(toUpdate, curTitle)
	}
	e.Shape().IterateBodies(f.GetBody(), args.ContentBlocks(), nil, func(bId content.BlockId, uId content.UnitId, curU, newU *content.Unit) {
		if newU != nil && newU.Content != "" {
			if curU == nil || curU.Content != newU.Content {
				toUpdate = append(toUpdate, newU)
			}
		}
	}, nil)

	// set flavor's full id
	for _, u := range toUpdate {
		u.ExtractId = f.ExtractId
		u.FlavorId = f.Id
		u.Language = f.Language
		u.FlavorType = f.Type
	}

	err = w.Content.InsertOrUpdateUnits(context.User, toUpdate)
	if err != nil {
		return nil, err
	}

	session.ClearDefaults()
	context.Query.Del("a")
	return nil, redirectToOther("/extract/" + e.UrlSlug + "/" + string(langCodeA) + "?" + context.Query.Encode())
}

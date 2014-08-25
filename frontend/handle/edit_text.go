package handle

import (
	"github.com/polyglottis/platform/content"
	"github.com/polyglottis/platform/frontend"
	"github.com/polyglottis/platform/language"
)

func (w *Worker) EditText(context *frontend.Context) ([]byte, error) {
	extract, langCodeA, langCodeB, err := w.readExtractAndLanguages(context)
	if err != nil {
		return nil, err
	}

	if fByTypeA, ok := extract.Flavors[langCodeA]; ok {
		if textA, ok := fByTypeA[content.Text]; ok {
			var textB *content.Flavor
			if len(langCodeB) != 0 {
				if fByTypeB, ok := extract.Flavors[langCodeB]; ok {
					if tB, ok := fByTypeB[content.Text]; ok {
						textB = tB[0]
					}
				}
			}
			return w.Server.EditText(context, extract, textA[0], textB)
		}
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

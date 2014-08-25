package handle

import (
	"github.com/polyglottis/platform/content"
	"github.com/polyglottis/platform/frontend"
	"github.com/polyglottis/platform/language"
)

func (w *Worker) EditText(context *frontend.Context) ([]byte, error) {
	id := content.ExtractId(context.Query.Get("id"))
	langA := context.Query.Get("a")
	langB := context.Query.Get("b")
	if len(id) == 0 || len(langA) == 0 {
		return nil, content.ErrInvalidInput
	}

	extract, err := w.Content.GetExtract(id)
	if err != nil {
		return nil, err
	}

	langCodeA, err := w.Language.GetCode(langA)
	if err != nil {
		return nil, err
	}

	var langCodeB language.Code
	if len(langB) != 0 {
		langCodeB, err = w.Language.GetCode(langB)
		if err != nil {
			return nil, err
		}
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

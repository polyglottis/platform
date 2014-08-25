package handle

import (
	"strconv"

	"github.com/polyglottis/platform/content"
	"github.com/polyglottis/platform/frontend"
)

func (w *Worker) Extract(context *frontend.Context) ([]byte, error) {
	extract, err := w.readExtract(context)
	if err != nil {
		return nil, err
	}
	return w.extract(context, extract)
}

func (w *Worker) extract(context *frontend.Context, extract *content.Extract) ([]byte, error) {
	return w.Server.Flavor(context, extract, nil, nil)
}

func (w *Worker) Flavor(context *frontend.Context) ([]byte, error) {
	extract, err := w.readExtract(context)
	if err != nil {
		return nil, err
	}

	langCode, err := w.Language.GetCode(context.Vars["language"])
	if err != nil {
		// language not found, fall back to extract
		return w.extract(context, extract)
	}

	if fByType, ok := extract.Flavors[langCode]; ok {
		a := newFlavorTriple(fByType, context, "a")

		b := &frontend.FlavorTriple{}
		langB := context.Query.Get("b")
		langCodeB, err := w.Language.GetCode(langB)
		if err == nil {
			if fByTypeB, ok := extract.Flavors[langCodeB]; ok {
				b = newFlavorTriple(fByTypeB, context, "b")
			}
		}
		return w.Server.Flavor(context, extract, a, b)
	}
	// flavor not found, fall back to extract
	return w.extract(context, extract)
}

func newFlavorTriple(fByType content.FlavorByType, context *frontend.Context, which string) *frontend.FlavorTriple {
	a := &frontend.FlavorTriple{}
	for _, data := range []struct {
		flavors        []*content.Flavor
		key            string
		insertionPoint **content.Flavor
	}{{
		flavors:        fByType[content.Audio],
		key:            which + "a",
		insertionPoint: &a.Audio,
	}, {
		flavors:        fByType[content.Text],
		key:            which + "t",
		insertionPoint: &a.Text,
	}, {
		flavors:        fByType[content.Transcript],
		key:            which + "p",
		insertionPoint: &a.Transcript,
	}} {
		if len(data.flavors) != 0 {
			idx := 0
			if want, err := strconv.Atoi(context.Query.Get(data.key)); err == nil {
				flavorId := content.FlavorId(want)
				for i, f := range data.flavors {
					if f.Id == flavorId {
						idx = i
						break
					}
				}
			}
			*data.insertionPoint = data.flavors[idx]
		}
	}
	return a
}

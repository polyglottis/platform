package frontend

import (
	"log"

	"github.com/polyglottis/platform/content"
)

func (w *Worker) Extract(context *Context) ([]byte, error) {
	slug := context.Vars["slug"]
	if len(slug) == 0 {
		return nil, content.ErrNotFound
	}

	id, err := w.Content.GetExtractId(slug)
	if err != nil {
		return nil, err
	}
	return w.extractById(context, id)
}

func (w *Worker) extractById(context *Context, id content.ExtractId) ([]byte, error) {
	extract, err := w.Content.GetExtract(id)
	if err != nil {
		return nil, err
	}
	return w.extract(context, extract)
}

func (w *Worker) extract(context *Context, extract *content.Extract) ([]byte, error) {
	for _, fByType := range extract.Flavors { // not great: not even deterministic...
		a := newFlavorTriple(fByType)
		return w.Server.Flavor(context, extract, a, &FlavorTriple{})
	}
	log.Println("Weird: extract with no flavor:", extract.Id)
	return nil, content.ErrNotFound
}

func (w *Worker) Flavor(context *Context) ([]byte, error) {
	slug := context.Vars["slug"]
	lang := context.Vars["language"]
	if len(slug) == 0 {
		return nil, content.ErrNotFound
	}

	id, err := w.Content.GetExtractId(slug)
	if err != nil {
		return nil, err
	}

	langCode, err := w.Language.GetCode(lang)
	if err != nil {
		// language not found, fall back to extract
		return w.extractById(context, id)
	}

	extract, err := w.Content.GetExtract(id)
	if err != nil {
		return nil, err
	}

	if fByType, ok := extract.Flavors[langCode]; ok {
		a := newFlavorTriple(fByType)
		return w.Server.Flavor(context, extract, a, &FlavorTriple{})
	}
	// flavor not found, fall back to extract
	return w.extract(context, extract)
}

func newFlavorTriple(fByType content.FlavorByType) *FlavorTriple {
	a := &FlavorTriple{}
	if audio, ok := fByType[content.Audio]; ok {
		a.Audio = audio[0]
	}
	if text, ok := fByType[content.Text]; ok {
		a.Text = text[0]
	}
	if transcript, ok := fByType[content.Transcript]; ok {
		a.Transcript = transcript[0]
	}
	return a
}

package content

import (
	"github.com/polyglottis/platform/language"
)

func (e *Extract) GetFlavor(langCode language.Code, flavorType FlavorType, id FlavorId) (*Flavor, error) {
	if fByType, ok := e.Flavors[langCode]; ok {
		if flavors, ok := fByType[flavorType]; ok {
			for _, f := range flavors {
				if f.Id == id {
					return f, nil
				}
			}
		}
	}
	return nil, ErrNotFound
}

func (f *Flavor) GetBody() BlockSlice {
	if f == nil {
		return nil
	}
	if len(f.Blocks) != 0 && f.Blocks[0][0].BlockId == 1 {
		return f.Blocks[1:]
	}
	return f.Blocks
}

func (f *Flavor) GetTitle() *Unit {
	if f == nil {
		return nil
	}
	if len(f.Blocks) != 0 && f.Blocks[0][0].BlockId == 1 {
		return f.Blocks[0][0]
	}
	return nil
}

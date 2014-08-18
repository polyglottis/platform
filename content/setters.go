package content

import (
	"github.com/polyglottis/platform/language"
)

func (e *Extract) SetId(id ExtractId) {
	e.Id = id
	for _, flavorTypes := range e.Flavors {
		for _, flavors := range flavorTypes {
			for _, flavor := range flavors {
				flavor.SetExtractId(id)
			}
		}
	}
}

func (f *Flavor) SetExtractId(id ExtractId) {
	f.ExtractId = id
	for _, block := range f.Blocks {
		for _, unit := range block {
			unit.ExtractId = id
		}
	}
}

func (f *Flavor) SetId(id FlavorId) {
	f.Id = id
	for _, block := range f.Blocks {
		for _, unit := range block {
			unit.FlavorId = id
		}
	}
}

func (e *Extract) SetFlavorLanguagesAndTypes() {
	for lang, fByType := range e.Flavors {
		for fType, flavors := range fByType {
			for _, f := range flavors {
				f.SetLanguage(lang)
				f.SetType(fType)
			}
		}
	}
}

func (f *Flavor) SetLanguage(lang language.Code) {
	f.Language = lang
	for _, block := range f.Blocks {
		for _, unit := range block {
			unit.Language = lang
		}
	}
}

func (f *Flavor) SetType(t FlavorType) {
	f.Type = t
	for _, block := range f.Blocks {
		for _, unit := range block {
			unit.FlavorType = t
		}
	}
}

package content

func (e *Extract) SetId(id ExtractId) {
	e.Id = id
	for _, flavor := range e.Flavors {
		flavor.SetExtractId(id)
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

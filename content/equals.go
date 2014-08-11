package content

func (f *Flavor) Equals(g *Flavor) bool {
	if f == nil {
		return g == nil
	} else if g == nil {
		return false
	}
	// now f != nil and g != nil
	if f.ExtractId != g.ExtractId ||
		f.Id != g.Id ||
		f.Summary != g.Summary ||
		f.Language != g.Language ||
		f.LanguageComment != g.LanguageComment ||
		f.Type != g.Type ||
		len(f.Blocks) != len(g.Blocks) {
		return false
	}

	for i, block := range f.Blocks {
		if len(block) != len(g.Blocks[i]) {
			return false
		}
		for j, u := range block {
			v := g.Blocks[i][j]
			if *u != *v {
				return false
			}
		}
	}
	return true
}

func (m *Metadata) Equals(n *Metadata) bool {
	if m == nil {
		return n == nil
	} else if n == nil {
		return false
	}
	return *m == *n
}

func (e *Extract) Equals(f *Extract) bool {
	if e == nil {
		return f == nil
	} else if f == nil {
		return false
	}
	// now e != nil and f != nil
	if e.Id != f.Id ||
		e.Type != f.Type ||
		e.UrlSlug != f.UrlSlug ||
		!e.Metadata.Equals(f.Metadata) ||
		len(e.Flavors) != len(f.Flavors) {
		return false
	}

	for i, flavor := range e.Flavors {
		if !flavor.Equals(f.Flavors[i]) {
			return false
		}
	}
	return true
}

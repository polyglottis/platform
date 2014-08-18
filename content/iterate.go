package content

func (s ExtractShape) IterateFlavorBody(f *Flavor, beforeBlock func(BlockId), unit func(BlockId, UnitId, *Unit), afterBlock func(BlockId)) {
	if len(f.Blocks) != 0 && f.Blocks[0][0].BlockId == 1 {
		// skip title
		s.IterateBody(f.Blocks[1:], beforeBlock, unit, afterBlock)
	} else {
		s.IterateBody(f.Blocks, beforeBlock, unit, afterBlock)
	}
}

func (s ExtractShape) IterateBody(body BlockSlice, beforeBlock func(BlockId), unit func(BlockId, UnitId, *Unit), afterBlock func(BlockId)) {
	nextBlockIdx := 0
	var nextBlock UnitSlice
	if nextBlockIdx < len(body) {
		nextBlock = body[nextBlockIdx]
	}
	for i, size := range s {
		blockId := BlockId(i + 2)
		if beforeBlock != nil {
			beforeBlock(blockId)
		}
		if nextBlock != nil && nextBlock[0].BlockId == blockId {
			if unit != nil {
				s.iterateUnits(size, blockId, nextBlock, unit)
			}
			nextBlockIdx++
			if nextBlockIdx < len(body) {
				nextBlock = body[nextBlockIdx]
			} else {
				nextBlock = nil
			}
		} else {
			if unit != nil {
				for j := 0; j < size; j++ {
					unit(blockId, UnitId(j+1), nil)
				}
			}
		}
		if afterBlock != nil {
			afterBlock(blockId)
		}
	}
}

func (s ExtractShape) iterateUnits(size int, blockId BlockId, units UnitSlice, f func(BlockId, UnitId, *Unit)) {
	nextUnitIdx := 0
	var nextUnit *Unit
	if nextUnitIdx < len(units) {
		nextUnit = units[nextUnitIdx]
	}
	for i := 0; i < size; i++ {
		unitId := UnitId(i + 1)
		if nextUnit != nil && nextUnit.Id == unitId {
			f(blockId, unitId, nextUnit)
			nextUnitIdx++
			if nextUnitIdx < len(units) {
				nextUnit = units[nextUnitIdx]
			} else {
				nextUnit = nil
			}
		} else {
			f(blockId, unitId, nil)
		}
	}
}

func (s ExtractShape) Union(f *Flavor) ExtractShape {
	if len(f.Blocks) == 0 {
		return s
	}
	lastBlockId := int(f.Blocks[len(f.Blocks)-1][0].BlockId)
	if lastBlockId > len(s)+1 {
		for i := len(s) + 1; i < lastBlockId; i++ {
			s = append(s, 0)
		}
	}
	for _, b := range f.Blocks {
		i := int(b[0].BlockId) - 2
		if i < 0 {
			continue
		}
		lastUnitId := int(b[len(b)-1].Id)
		if s[i] < lastUnitId {
			s[i] = lastUnitId
		}
	}
	return s
}

func (s ExtractShape) Equals(t ExtractShape) bool {
	if len(s) != len(t) {
		return false
	}
	for i, si := range s {
		if si != t[i] {
			return false
		}
	}
	return true
}

func (e *Extract) Shape() ExtractShape {
	s := ExtractShape{}
	for _, fByType := range e.Flavors {
		for _, flavors := range fByType {
			for _, f := range flavors {
				s = s.Union(f)
			}
		}
	}
	return s
}

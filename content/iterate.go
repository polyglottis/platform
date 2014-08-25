package content

func (s ExtractShape) IterateFlavorBody(f *Flavor, beforeBlock func(BlockId), unit func(BlockId, UnitId, *Unit), afterBlock func(BlockId)) {
	if len(f.Blocks) != 0 && f.Blocks[0][0].BlockId == 1 {
		// skip title
		s.IterateBody(f.Blocks[1:], beforeBlock, unit, afterBlock)
	} else {
		s.IterateBody(f.Blocks, beforeBlock, unit, afterBlock)
	}
}
func (s ExtractShape) IterateFlavorBodies(fA, fB *Flavor, beforeBlock func(BlockId), unit func(BlockId, UnitId, *Unit, *Unit), afterBlock func(BlockId)) {
	var blocksA, blocksB BlockSlice
	if fA != nil {
		blocksA = fA.Blocks
	}
	if fB != nil {
		blocksB = fB.Blocks
	}
	if len(blocksA) != 0 && blocksA[0][0].BlockId == 1 {
		blocksA = blocksA[1:]
	}
	if len(blocksB) != 0 && blocksB[0][0].BlockId == 1 {
		blocksB = blocksB[1:]
	}
	s.IterateBodies(blocksA, blocksB, beforeBlock, unit, afterBlock)
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
func (s ExtractShape) IterateBodies(bodyA, bodyB BlockSlice, beforeBlock func(BlockId), unit func(BlockId, UnitId, *Unit, *Unit), afterBlock func(BlockId)) {
	nextBlockIdxA, nextBlockIdxB := 0, 0
	var nextBlockA, nextBlockB UnitSlice
	if nextBlockIdxA < len(bodyA) {
		nextBlockA = bodyA[nextBlockIdxA]
	}
	if nextBlockIdxB < len(bodyB) {
		nextBlockB = bodyB[nextBlockIdxB]
	}
	for i, size := range s {
		blockId := BlockId(i + 2)
		if beforeBlock != nil {
			beforeBlock(blockId)
		}
		var curBlockA, curBlockB UnitSlice
		if nextBlockA != nil && nextBlockA[0].BlockId == blockId {
			curBlockA = nextBlockA
			nextBlockIdxA++
			if nextBlockIdxA < len(bodyA) {
				nextBlockA = bodyA[nextBlockIdxA]
			} else {
				nextBlockA = nil
			}
		}
		if nextBlockB != nil && nextBlockB[0].BlockId == blockId {
			curBlockB = nextBlockB
			nextBlockIdxB++
			if nextBlockIdxB < len(bodyB) {
				nextBlockB = bodyB[nextBlockIdxB]
			} else {
				nextBlockB = nil
			}
		}
		if unit != nil {
			if curBlockA == nil && curBlockB == nil {
				for j := 0; j < size; j++ {
					unit(blockId, UnitId(j+1), nil, nil)
				}
			} else {
				s.iterateUnitsAB(size, blockId, curBlockA, curBlockB, unit)
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

func (s ExtractShape) iterateUnitsAB(size int, blockId BlockId, unitsA, unitsB UnitSlice, f func(BlockId, UnitId, *Unit, *Unit)) {
	nextUnitIdxA, nextUnitIdxB := 0, 0
	var nextUnitA, nextUnitB *Unit
	if nextUnitIdxA < len(unitsA) {
		nextUnitA = unitsA[nextUnitIdxA]
	}
	if nextUnitIdxB < len(unitsB) {
		nextUnitB = unitsB[nextUnitIdxB]
	}
	for i := 0; i < size; i++ {
		unitId := UnitId(i + 1)
		var curUnitA, curUnitB *Unit
		if nextUnitA != nil && nextUnitA.Id == unitId {
			curUnitA = nextUnitA
			nextUnitIdxA++
			if nextUnitIdxA < len(unitsA) {
				nextUnitA = unitsA[nextUnitIdxA]
			} else {
				nextUnitA = nil
			}
		}
		if nextUnitB != nil && nextUnitB.Id == unitId {
			curUnitB = nextUnitB
			nextUnitIdxB++
			if nextUnitIdxB < len(unitsB) {
				nextUnitB = unitsB[nextUnitIdxB]
			} else {
				nextUnitB = nil
			}
		}
		f(blockId, unitId, curUnitA, curUnitB)
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

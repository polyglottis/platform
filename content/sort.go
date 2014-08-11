package content

func (bs BlockSlice) Len() int      { return len(bs) }
func (bs BlockSlice) Swap(i, j int) { bs[i], bs[j] = bs[j], bs[i] }
func (bs BlockSlice) Less(i, j int) bool {
	switch {
	case len(bs[i]) == 0:
		return true
	case len(bs[j]) == 0:
		return false
	default:
		return bs[i][0].BlockId < bs[j][0].BlockId
	}
}

func (us UnitSlice) Len() int           { return len(us) }
func (us UnitSlice) Swap(i, j int)      { us[i], us[j] = us[j], us[i] }
func (us UnitSlice) Less(i, j int) bool { return us[i].Id < us[j].Id }

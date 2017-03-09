package ribbon

type Helix struct {
	ChainID    string
	InitSeqNum int
	EndSeqNum  int
}

func (h *Helix) Contains(r *Residue) bool {
	return r.ChainID == h.ChainID &&
		r.ResSeq >= h.InitSeqNum && r.ResSeq <= h.EndSeqNum
}

package ribbon

type Strand struct {
	ChainID    string
	InitSeqNum int
	EndSeqNum  int
}

func (s *Strand) Contains(r *Residue) bool {
	return r.ChainID == s.ChainID &&
		r.ResSeq >= s.InitSeqNum && r.ResSeq <= s.EndSeqNum
}

package ribbon

type ResidueType int

const (
	_ ResidueType = iota
	ResidueTypeOther
	ResidueTypeHelix
	ResidueTypeStrand
)

type Residue struct {
	Type    ResidueType
	ResSeq  int
	ChainID string
	Atoms   map[string]*Atom
}

func NewResidue(atoms []*Atom) *Residue {
	resSeq := atoms[0].ResSeq
	chainID := atoms[0].ChainID
	m := make(map[string]*Atom)
	for _, a := range atoms {
		m[a.Name] = a
	}
	return &Residue{ResidueTypeOther, resSeq, chainID, m}
}

func ResiduesForAtoms(atoms []*Atom) []*Residue {
	var residues []*Residue
	var group []*Atom
	previous := -1
	for _, atom := range atoms {
		value := atom.ResSeq
		if value != previous && group != nil {
			residues = append(residues, NewResidue(group))
			group = nil
		}
		group = append(group, atom)
		previous = value
	}
	residues = append(residues, NewResidue(group))
	return residues
}

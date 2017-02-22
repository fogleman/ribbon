package ribbon

type Residue struct {
	Number int
	Name   string
	Chain  string
	Atoms  map[string]*Atom
}

func NewResidue(atoms []*Atom) *Residue {
	number := atoms[0].ResSeq
	name := atoms[0].ResName
	chain := atoms[0].ChainID
	m := make(map[string]*Atom)
	for _, a := range atoms {
		m[a.Name] = a
	}
	return &Residue{number, name, chain, m}
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

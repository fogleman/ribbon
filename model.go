package ribbon

type Model struct {
	Atoms    []*Atom
	Residues []*Residue
}

func NewModel(atoms []*Atom) *Model {
	residues := ResiduesForAtoms(atoms)
	return &Model{atoms, residues}
}

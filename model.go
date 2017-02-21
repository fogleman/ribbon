package ribbon

type Model struct {
	Atoms        []*Atom
	Residues     []*Residue
	Polypeptides []*Polypeptide
}

func NewModel(atoms []*Atom) *Model {
	residues := ResiduesForAtoms(atoms)
	polypeptides := PolypeptidesForResidues(residues)
	return &Model{atoms, residues, polypeptides}
}

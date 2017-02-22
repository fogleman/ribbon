package ribbon

type Model struct {
	Atoms        []*Atom
	Helixes      []*Helix
	Strands      []*Strand
	Residues     []*Residue
	Polypeptides []*Polypeptide
}

func NewModel(atoms []*Atom, helixes []*Helix, strands []*Strand) *Model {
	residues := ResiduesForAtoms(atoms)
	polypeptides := PolypeptidesForResidues(residues)
	return &Model{atoms, helixes, strands, residues, polypeptides}
}

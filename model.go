package ribbon

type Model struct {
	Atoms    []*Atom
	Helixes  []*Helix
	Strands  []*Strand
	Residues []*Residue
	Chains   []*Chain
}

func NewModel(atoms []*Atom, helixes []*Helix, strands []*Strand) *Model {
	residues := ResiduesForAtoms(atoms)
	chains := ChainsForResidues(residues)
	return &Model{atoms, helixes, strands, residues, chains}
}

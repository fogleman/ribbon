package pdb

type Model struct {
	Atoms       []*Atom
	HetAtoms    []*Atom
	Connections []*Connection
	Helixes     []*Helix
	Strands     []*Strand
	BioMatrixes []Matrix
	SymMatrixes []Matrix
	Residues    []*Residue
	Chains      []*Chain
}

func (model *Model) RemoveChain(chainID string) {
	atoms := model.Atoms[:0]
	for _, atom := range model.Atoms {
		if atom.ChainID != chainID {
			atoms = append(atoms, atom)
		}
	}
	model.Atoms = atoms

	hetAtoms := model.HetAtoms[:0]
	for _, atom := range model.HetAtoms {
		if atom.ChainID != chainID {
			hetAtoms = append(hetAtoms, atom)
		}
	}
	model.HetAtoms = hetAtoms

	residues := model.Residues[:0]
	for _, residue := range model.Residues {
		if residue.ChainID != chainID {
			residues = append(residues, residue)
		}
	}
	model.Residues = residues

	chains := model.Chains[:0]
	for _, chain := range model.Chains {
		if chain.ChainID != chainID {
			chains = append(chains, chain)
		}
	}
	model.Chains = chains
}

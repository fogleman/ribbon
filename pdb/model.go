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

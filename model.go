package ribbon

import "github.com/fogleman/fauxgl"

type Model struct {
	Atoms              []*Atom
	Helixes            []*Helix
	Strands            []*Strand
	Residues           []*Residue
	Chains             []*Chain
	BiologicalMatrixes []fauxgl.Matrix
	SymmetryMatrixes   []fauxgl.Matrix
}

func NewModel(atoms []*Atom, helixes []*Helix, strands []*Strand) *Model {
	residues := ResiduesForAtoms(atoms)
	chains := ChainsForResidues(residues)
	for _, r := range residues {
		for _, h := range helixes {
			if r.Chain == h.Chain && r.Number >= h.Start && r.Number <= h.End {
				r.Type = ResidueTypeHelix
			}
		}
		for _, s := range strands {
			if r.Chain == s.Chain && r.Number >= s.Start && r.Number <= s.End {
				r.Type = ResidueTypeStrand
			}
		}
	}
	return &Model{atoms, helixes, strands, residues, chains, nil, nil}
}

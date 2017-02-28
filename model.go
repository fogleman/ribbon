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
			if r.ChainID == h.ChainID && r.ResSeq >= h.InitSeqNum && r.ResSeq <= h.EndSeqNum {
				r.Type = ResidueTypeHelix
			}
		}
		for _, s := range strands {
			if r.ChainID == s.ChainID && r.ResSeq >= s.InitSeqNum && r.ResSeq <= s.EndSeqNum {
				r.Type = ResidueTypeStrand
			}
		}
	}
	return &Model{atoms, helixes, strands, residues, chains, nil, nil}
}

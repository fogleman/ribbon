package pdb

type Residue struct {
	ResName     string
	ChainID     string
	ResSeq      int
	Atoms       []*Atom
	AtomsByName map[string]*Atom
	Secondary   Secondary
}

func newResidue(atoms []*Atom) *Residue {
	m := make(map[string]*Atom)
	for _, a := range atoms {
		m[a.Name] = a
	}
	residue := Residue{}
	residue.ResName = atoms[0].ResName
	residue.ChainID = atoms[0].ChainID
	residue.ResSeq = atoms[0].ResSeq
	residue.Atoms = atoms
	residue.AtomsByName = m
	residue.Secondary = SecondaryCoil
	return &residue
}

func residuesForAtoms(atoms []*Atom, helixes []*Helix, strands []*Strand) []*Residue {
	var residues []*Residue
	var group []*Atom
	previous := -1
	for _, atom := range atoms {
		value := atom.ResSeq
		if value != previous && group != nil {
			residues = append(residues, newResidue(group))
			group = nil
		}
		group = append(group, atom)
		previous = value
	}
	residues = append(residues, newResidue(group))
	for _, r := range residues {
		for _, h := range helixes {
			if r.ChainID == h.InitChainID && r.ResSeq >= h.InitSeqNum && r.ResSeq <= h.EndSeqNum {
				r.Secondary = SecondaryHelix
			}
		}
		for _, s := range strands {
			if r.ChainID == s.InitChainID && r.ResSeq >= s.InitSeqNum && r.ResSeq <= s.EndSeqNum {
				r.Secondary = SecondaryStrand
			}
		}
	}
	return residues
}

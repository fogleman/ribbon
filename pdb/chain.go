package pdb

type Chain struct {
	ChainID  string
	Residues []*Residue
}

func newChain(residues []*Residue) *Chain {
	chain := Chain{}
	chain.ChainID = residues[0].ChainID
	chain.Residues = residues
	return &chain
}

func chainsForResidues(residues []*Residue) []*Chain {
	var chains []*Chain
	var group []*Residue
	previous := residues[0]
	for _, residue := range residues {
		distance := residue.distance(previous)
		if residue.ChainID != previous.ChainID || distance > 4 {
			if group != nil {
				chains = append(chains, newChain(group))
				group = nil
			}
		}
		group = append(group, residue)
		previous = residue
	}
	if group != nil {
		chains = append(chains, newChain(group))
	}
	return chains
}

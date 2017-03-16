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
	previous := ""
	for _, residue := range residues {
		value := residue.ChainID
		if value != previous && group != nil {
			chains = append(chains, newChain(group))
			group = nil
		}
		group = append(group, residue)
		previous = value
	}
	if group != nil {
		chains = append(chains, newChain(group))
	}
	return chains
}

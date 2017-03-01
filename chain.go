package ribbon

import "github.com/fogleman/fauxgl"

type Chain struct {
	PeptidePlanes []*PeptidePlane
}

func NewChain(planes []*PeptidePlane) *Chain {
	var previous fauxgl.Vector
	for i, p := range planes {
		if i > 0 && p.Side.Dot(previous) < 0 {
			p.Side = p.Side.Negate()
			p.Normal = p.Normal.Negate()
		}
		previous = p.Side
	}
	return &Chain{planes}
}

func ChainsForResidues(residues []*Residue) []*Chain {
	var chains []*Chain
	var planes []*PeptidePlane
	for i := 0; i < len(residues)-2; i++ {
		r1 := residues[i]
		r2 := residues[i+1]
		r3 := residues[i+2]
		p := NewPeptidePlane(r1, r2, r3)
		if p != nil {
			planes = append(planes, p)
		} else if planes != nil {
			chains = append(chains, NewChain(planes))
			planes = nil
		}
	}
	if planes != nil {
		chains = append(chains, NewChain(planes))
	}
	return chains
}

func (c *Chain) Mesh() *fauxgl.Mesh {
	mesh := fauxgl.NewEmptyMesh()
	for i := 0; i < len(c.PeptidePlanes)-3; i++ {
		// TODO: handle ends better
		pp1 := c.PeptidePlanes[i]
		pp2 := c.PeptidePlanes[i+1]
		pp3 := c.PeptidePlanes[i+2]
		pp4 := c.PeptidePlanes[i+3]
		m := createSegmentMesh(pp1, pp2, pp3, pp4)
		mesh.Add(m)
	}
	return mesh
}

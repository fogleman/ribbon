package ribbon

import "github.com/fogleman/fauxgl"

type Chain struct {
	PeptidePlanes []*PeptidePlane
}

func NewChain(planes []*PeptidePlane) *Chain {
	// var previous fauxgl.Vector
	// for i, p := range planes {
	// 	if i > 0 && p.Side.Dot(previous) < 0 {
	// 		p.Side = p.Side.Negate()
	// 		// p.Normal = p.Normal.Negate()
	// 	}
	// 	previous = p.Side
	// }
	return &Chain{planes}
}

func ChainsForResidues(residues []*Residue) []*Chain {
	var chains []*Chain
	var planes []*PeptidePlane
	for i := 0; i < len(residues)-1; i++ {
		r1 := residues[i]
		r2 := residues[i+1]
		p := NewPeptidePlane(r1, r2)
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

func (c *Chain) Ribbon(width, height float64) *fauxgl.Mesh {
	const n = 64
	var triangles []*fauxgl.Triangle
	var lines []*fauxgl.Line
	for i := 0; i < len(c.PeptidePlanes)-3; i++ {
		// TODO: handle ends
		// TODO: handle plane flips
		p0 := c.PeptidePlanes[i]
		p1 := c.PeptidePlanes[i+1]
		p2 := c.PeptidePlanes[i+2]
		p3 := c.PeptidePlanes[i+3]
		var splines [2][2][]fauxgl.Vector
		for u := 0; u < 2; u++ {
			for v := 0; v < 2; v++ {
				w := float64(u*2-1) * width / 2
				h := float64(v*2-1) * height / 2
				g0 := p0.Position.Add(p0.Side.MulScalar(w)).Add(p0.Normal.MulScalar(h + 1.5))
				g1 := p1.Position.Add(p1.Side.MulScalar(w)).Add(p1.Normal.MulScalar(h + 1.5))
				g2 := p2.Position.Add(p2.Side.MulScalar(w)).Add(p2.Normal.MulScalar(h + 1.5))
				g3 := p3.Position.Add(p3.Side.MulScalar(w)).Add(p3.Normal.MulScalar(h + 1.5))
				splines[u][v] = Spline(g0, g1, g2, g3, n)
			}
		}
		for j := 0; j < n; j++ {
			p000 := splines[0][0][j]
			p001 := splines[0][0][j+1]
			p010 := splines[0][1][j]
			p011 := splines[0][1][j+1]
			p100 := splines[1][0][j]
			p101 := splines[1][0][j+1]
			p110 := splines[1][1][j]
			p111 := splines[1][1][j+1]
			triangles = triangulateQuad(triangles, p000, p100, p101, p001)
			triangles = triangulateQuad(triangles, p011, p111, p110, p010)
			triangles = triangulateQuad(triangles, p110, p111, p101, p100)
			triangles = triangulateQuad(triangles, p000, p001, p011, p010)
			lines = append(lines, fauxgl.NewLineForPoints(p000, p001))
			lines = append(lines, fauxgl.NewLineForPoints(p010, p011))
			lines = append(lines, fauxgl.NewLineForPoints(p100, p101))
			lines = append(lines, fauxgl.NewLineForPoints(p110, p111))
		}
	}
	return fauxgl.NewMesh(triangles, lines)
}

func triangulateQuad(triangles []*fauxgl.Triangle, p1, p2, p3, p4 fauxgl.Vector) []*fauxgl.Triangle {
	triangles = append(triangles, fauxgl.NewTriangleForPoints(p1, p2, p3))
	triangles = append(triangles, fauxgl.NewTriangleForPoints(p1, p3, p4))
	return triangles
}

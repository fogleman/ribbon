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
		r1 := p1.Residue1
		r2 := p2.Residue1
		var splines1 [2][2][]fauxgl.Vector
		var splines2 [2][2][]fauxgl.Vector
		for u := 0; u < 2; u++ {
			for v := 0; v < 2; v++ {
				w1 := float64(u*2-1) * 0.25
				h1 := float64(v*2-1) * 0.25
				if r1.Type == ResidueTypeHelix {
					w1 = float64(u*2-1) * width / 2
					h1 = float64(v*2-1) * height / 16
					h1 += 1.5
				} else if r1.Type == ResidueTypeStrand {
					if r2.Type == ResidueTypeStrand {
						w1 = float64(u*2-1) * width / 2
						h1 = float64(v*2-1) * height / 2
					} else {
						w1 = float64(u*2-1) * width
						h1 = float64(v*2-1) * height / 2
					}
				}
				w2 := float64(u*2-1) * 0.25
				h2 := float64(v*2-1) * 0.25
				if r2.Type == ResidueTypeHelix {
					w2 = float64(u*2-1) * width / 2
					h2 = float64(v*2-1) * height / 16
					h2 += 1.5
				} else if r2.Type == ResidueTypeStrand {
					w2 = float64(u*2-1) * width / 2
					h2 = float64(v*2-1) * height / 2
				}
				splines1[u][v] = SplineForPlanes(p0, p1, p2, p3, n, w1, h1)
				splines2[u][v] = SplineForPlanes(p0, p1, p2, p3, n, w2, h2)
			}
		}
		for j := 0; j < n; j++ {
			t1 := float64(j) / float64(n)
			t2 := float64(j+1) / float64(n)
			if r2.Type == ResidueTypeStrand && r1.Type != ResidueTypeStrand {
				if t1 < 0.5 {
					t1 = 0
				} else {
					t1 = 1
				}
				if t2 < 0.5 {
					t2 = 0
				} else {
					t2 = 1
				}
			}
			p000 := splines1[0][0][j].Lerp(splines2[0][0][j], t1)
			p001 := splines1[0][0][j+1].Lerp(splines2[0][0][j+1], t2)
			p010 := splines1[0][1][j].Lerp(splines2[0][1][j], t1)
			p011 := splines1[0][1][j+1].Lerp(splines2[0][1][j+1], t2)
			p100 := splines1[1][0][j].Lerp(splines2[1][0][j], t1)
			p101 := splines1[1][0][j+1].Lerp(splines2[1][0][j+1], t2)
			p110 := splines1[1][1][j].Lerp(splines2[1][1][j], t1)
			p111 := splines1[1][1][j+1].Lerp(splines2[1][1][j+1], t2)
			triangles = triangulateQuad(triangles, p000, p100, p101, p001)
			triangles = triangulateQuad(triangles, p011, p111, p110, p010)
			triangles = triangulateQuad(triangles, p110, p111, p101, p100)
			triangles = triangulateQuad(triangles, p000, p001, p011, p010)
			// if r.Type == ResidueTypeHelix {
			lines = append(lines, fauxgl.NewLineForPoints(p000, p001))
			lines = append(lines, fauxgl.NewLineForPoints(p010, p011))
			lines = append(lines, fauxgl.NewLineForPoints(p100, p101))
			lines = append(lines, fauxgl.NewLineForPoints(p110, p111))
			// }
		}
	}
	return fauxgl.NewMesh(triangles, lines)
}

func triangulateQuad(triangles []*fauxgl.Triangle, p1, p2, p3, p4 fauxgl.Vector) []*fauxgl.Triangle {
	triangles = append(triangles, fauxgl.NewTriangleForPoints(p1, p2, p3))
	triangles = append(triangles, fauxgl.NewTriangleForPoints(p1, p3, p4))
	return triangles
}

package ribbon

import "github.com/fogleman/fauxgl"

type Polypeptide struct {
	PeptidePlanes []*PeptidePlane
}

func NewPolypeptide(planes []*PeptidePlane) *Polypeptide {
	return &Polypeptide{planes}
}

func PolypeptidesForResidues(residues []*Residue) []*Polypeptide {
	var polypeptides []*Polypeptide
	var planes []*PeptidePlane
	for i := 0; i < len(residues)-1; i++ {
		r1 := residues[i]
		r2 := residues[i+1]
		p := NewPeptidePlane(r1, r2)
		if p != nil {
			planes = append(planes, p)
		} else if planes != nil {
			polypeptides = append(polypeptides, NewPolypeptide(planes))
			planes = nil
		}
	}
	if planes != nil {
		polypeptides = append(polypeptides, NewPolypeptide(planes))
	}
	return polypeptides
}

func (pp *Polypeptide) Ribbon(width, height float64) *fauxgl.Mesh {
	const n = 64
	var triangles []*fauxgl.Triangle
	for i := 0; i < len(pp.PeptidePlanes)-3; i++ {
		// TODO: handle ends
		// TODO: handle plane flips
		p0 := pp.PeptidePlanes[i]
		p1 := pp.PeptidePlanes[i+1]
		p2 := pp.PeptidePlanes[i+2]
		p3 := pp.PeptidePlanes[i+3]
		var splines [2][2][]fauxgl.Vector
		for u := 0; u < 2; u++ {
			for v := 0; v < 2; v++ {
				w := float64(u*2-1) * width / 2
				h := float64(v*2-1) * height / 2
				g0 := p0.Position.Add(p0.Side.MulScalar(w)).Add(p0.Normal.MulScalar(h))
				g1 := p1.Position.Add(p1.Side.MulScalar(w)).Add(p1.Normal.MulScalar(h))
				g2 := p2.Position.Add(p2.Side.MulScalar(w)).Add(p2.Normal.MulScalar(h))
				g3 := p3.Position.Add(p3.Side.MulScalar(w)).Add(p3.Normal.MulScalar(h))
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
		}
	}
	return fauxgl.NewTriangleMesh(triangles)
}

func triangulateQuad(triangles []*fauxgl.Triangle, p1, p2, p3, p4 fauxgl.Vector) []*fauxgl.Triangle {
	triangles = append(triangles, fauxgl.NewTriangleForPoints(p1, p2, p3))
	triangles = append(triangles, fauxgl.NewTriangleForPoints(p1, p3, p4))
	return triangles
}

func makeSegment(p0, p1 fauxgl.Vector, r float64) *fauxgl.Mesh {
	p := p0.Add(p1).MulScalar(0.5)
	h := p0.Distance(p1) * 2
	up := p1.Sub(p0).Normalize()
	mesh := fauxgl.NewCylinder(30, false)
	mesh.Transform(fauxgl.Orient(p, fauxgl.V(r, r, h), up, 0))
	return mesh
}

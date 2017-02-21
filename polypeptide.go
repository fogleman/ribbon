package ribbon

import (
	"math/rand"

	"github.com/fogleman/fauxgl"
)

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

func (pp *Polypeptide) Ribbon(width, thickness float64) *fauxgl.Mesh {
	// var lines []*fauxgl.Line
	mesh := fauxgl.NewEmptyMesh()
	for i := 0; i < len(pp.PeptidePlanes)-3; i++ {
		// TODO: handle ends
		p0 := pp.PeptidePlanes[i].Position
		p1 := pp.PeptidePlanes[i+1].Position
		p2 := pp.PeptidePlanes[i+2].Position
		p3 := pp.PeptidePlanes[i+3].Position
		points := Spline(p0, p1, p2, p3, 8)
		for j := 0; j < len(points)-1; j++ {
			s := makeSegment(points[j], points[j+1], 1)
			mesh.Add(s)
			// line := fauxgl.NewLineForPoints(points[j], points[j+1])
			// lines = append(lines, line)
		}
	}
	c := fauxgl.Color{rand.Float64(), rand.Float64(), rand.Float64(), 1}
	for _, t := range mesh.Triangles {
		t.V1.Color = c
		t.V2.Color = c
		t.V3.Color = c
	}
	return mesh
	// return fauxgl.NewLineMesh(lines)
}

func makeSegment(p0, p1 fauxgl.Vector, r float64) *fauxgl.Mesh {
	p := p0.Add(p1).MulScalar(0.5)
	h := p0.Distance(p1) * 2
	up := p1.Sub(p0).Normalize()
	mesh := fauxgl.NewCylinder(30, false)
	mesh.Transform(fauxgl.Orient(p, fauxgl.V(r, r, h), up, 0))
	return mesh
}

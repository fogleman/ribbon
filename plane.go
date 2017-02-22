package ribbon

import "github.com/fogleman/fauxgl"

type PeptidePlane struct {
	Position fauxgl.Vector
	Normal   fauxgl.Vector
	Forward  fauxgl.Vector
	Side     fauxgl.Vector
}

func NewPeptidePlane(r1, r2 *Residue) *PeptidePlane {
	if r1.Chain != r2.Chain {
		return nil
	}
	ca1 := r1.Atoms["CA"]
	ca2 := r2.Atoms["CA"]
	o1 := r1.Atoms["O"]
	if ca1 == nil || ca2 == nil || o1 == nil {
		return nil
	}
	a := ca2.Position.Sub(ca1.Position)
	b := o1.Position.Sub(ca1.Position)
	c := a.Cross(b)
	d := c.Cross(a)
	p := ca1.Position.Add(ca2.Position).DivScalar(2)
	return &PeptidePlane{p, c.Normalize(), a.Normalize(), d.Normalize()}
}

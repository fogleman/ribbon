package ribbon

import "github.com/fogleman/fauxgl"

type PeptidePlane struct {
	Residue1 *Residue
	Residue2 *Residue
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
	p := ca1.Position //.Add(ca2.Position).DivScalar(2)
	return &PeptidePlane{r1, r2, p, c.Normalize(), a.Normalize(), d.Normalize()}
}

func (p *PeptidePlane) Point(v fauxgl.Vector) fauxgl.Vector {
	result := p.Position
	result = result.Add(p.Side.MulScalar(v.X))
	result = result.Add(p.Normal.MulScalar(v.Y))
	result = result.Add(p.Forward.MulScalar(v.Z))
	return result
}

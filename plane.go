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
	if r1.ChainID != r2.ChainID {
		return nil
	}
	ca1 := r1.Atoms["CA"]
	ca2 := r2.Atoms["CA"]
	o1 := r1.Atoms["O"]
	if ca1 == nil || ca2 == nil || o1 == nil {
		return nil
	}
	a := ca2.Position.Sub(ca1.Position).Normalize()
	b := o1.Position.Sub(ca1.Position).Normalize()
	c := a.Cross(b).Normalize()
	d := c.Cross(a).Normalize()
	p := ca1.Position.Add(ca2.Position).DivScalar(2)
	return &PeptidePlane{r1, r2, p, c, a, d}
}

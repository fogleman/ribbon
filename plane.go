package ribbon

import "github.com/fogleman/fauxgl"

type PeptidePlane struct {
	Residue1 *Residue
	Residue2 *Residue
	Residue3 *Residue
	Position fauxgl.Vector
	Normal   fauxgl.Vector
	Forward  fauxgl.Vector
	Side     fauxgl.Vector
	Flipped  bool
}

func NewPeptidePlane(r1, r2, r3 *Residue) *PeptidePlane {
	if r1.ChainID != r2.ChainID || r2.ChainID != r3.ChainID {
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
	return &PeptidePlane{r1, r2, r3, p, c, a, d, false}
}

func (pp *PeptidePlane) Transition() (s1, s2 Secondary) {
	t1 := pp.Residue1.Secondary
	t2 := pp.Residue2.Secondary
	t3 := pp.Residue3.Secondary
	s1 = t2
	s2 = t2
	if t2 > t1 && t2 == t3 {
		s1 = t1
	}
	if t2 > t3 && t1 == t2 {
		s2 = t3
	}
	return
}

func (pp *PeptidePlane) Flip() {
	pp.Side = pp.Side.Negate()
	pp.Normal = pp.Normal.Negate()
	pp.Flipped = !pp.Flipped
}

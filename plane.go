package ribbon

import (
	"github.com/fogleman/fauxgl"
	"github.com/fogleman/ribbon/pdb"
)

type PeptidePlane struct {
	Residue1 *pdb.Residue
	Residue2 *pdb.Residue
	Residue3 *pdb.Residue
	Position fauxgl.Vector
	Normal   fauxgl.Vector
	Forward  fauxgl.Vector
	Side     fauxgl.Vector
	Flipped  bool
}

func NewPeptidePlane(r1, r2, r3 *pdb.Residue) *PeptidePlane {
	// TODO: handle missing required atoms
	ca1 := atomPosition(r1.AtomsByName["CA"])
	ca2 := atomPosition(r2.AtomsByName["CA"])
	o1 := atomPosition(r1.AtomsByName["O"])
	a := ca2.Sub(ca1).Normalize()
	b := o1.Sub(ca1).Normalize()
	c := a.Cross(b).Normalize()
	d := c.Cross(a).Normalize()
	p := ca1.Add(ca2).DivScalar(2)
	return &PeptidePlane{r1, r2, r3, p, c, a, d, false}
}

func (pp *PeptidePlane) Transition() (s1, s2 pdb.Secondary) {
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

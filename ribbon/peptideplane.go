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
	if _, ok := r1.AtomsByName["CA"]; !ok {
		return nil
	}
	if _, ok := r2.AtomsByName["CA"]; !ok {
		return nil
	}
	if _, ok := r1.AtomsByName["O"]; !ok {
		return nil
	}
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

func (pp *PeptidePlane) Transition() (type1, type2 pdb.ResidueType) {
	t1 := pp.Residue1.Type
	t2 := pp.Residue2.Type
	t3 := pp.Residue3.Type
	type1 = t2
	type2 = t2
	if t2 > t1 && t2 == t3 {
		type1 = t1
	}
	if t2 > t3 && t1 == t2 {
		type2 = t3
	}
	return
}

func (pp *PeptidePlane) Flip() {
	pp.Side = pp.Side.Negate()
	pp.Normal = pp.Normal.Negate()
	pp.Flipped = !pp.Flipped
}

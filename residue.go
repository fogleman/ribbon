package ribbon

type Residue struct {
	Name  string
	Atoms map[string]*Atom
}

func NewResidue(atoms []*Atom) *Residue {
	name := atoms[0].ResName
	m := make(map[string]*Atom)
	for _, a := range atoms {
		m[a.Name] = a
	}
	return &Residue{name, m}
}

func ResiduesForAtoms(atoms []*Atom) []*Residue {
	var residues []*Residue
	var group []*Atom
	previous := -1
	for _, atom := range atoms {
		value := atom.ResSeq
		if value != previous && group != nil {
			residues = append(residues, NewResidue(group))
			group = nil
		}
		group = append(group, atom)
		previous = value
	}
	residues = append(residues, NewResidue(group))
	return residues
}

type ResiduePlane struct {
	Position Vector
	Normal   Vector
	Forward  Vector
	Side     Vector
}

func NewResiduePlane(r1, r2 *Residue) *ResiduePlane {
	ca1 := r1.Atoms["CA"].Position
	ca2 := r2.Atoms["CA"].Position
	o1 := r1.Atoms["O"].Position
	a := ca2.Sub(ca1)
	b := o1.Sub(ca1)
	c := a.Cross(b)
	d := c.Cross(a)
	p := ca1.Add(ca2).DivScalar(2)
	return &ResiduePlane{p, c.Normalize(), a.Normalize(), d.Normalize()}
}

func (a *ResiduePlane) Lerp(b *ResiduePlane, t float64) *ResiduePlane {
	p := ResiduePlane{}
	p.Position = a.Position.Lerp(b.Position, t)
	p.Normal = a.Normal.Lerp(b.Normal, t)
	p.Forward = a.Forward.Lerp(b.Forward, t)
	p.Side = a.Side.Lerp(b.Side, t)
	return &p
}

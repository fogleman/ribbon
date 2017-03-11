package ribbon

import "github.com/fogleman/fauxgl"

type Model struct {
	Atoms              []*Atom
	HetAtoms           []*Atom
	Connections        []Connection
	Helixes            []*Helix
	Strands            []*Strand
	Residues           []*Residue
	Chains             []*Chain
	BiologicalMatrixes []fauxgl.Matrix
	SymmetryMatrixes   []fauxgl.Matrix
}

func NewModel(atoms, hetAtoms []*Atom, connections []Connection, helixes []*Helix, strands []*Strand) *Model {
	residues := ResiduesForAtoms(atoms)
	chains := ChainsForResidues(residues)
	for _, r := range residues {
		for _, h := range helixes {
			if h.Contains(r) {
				r.Type = ResidueTypeHelix
			}
		}
		for _, s := range strands {
			if s.Contains(r) {
				r.Type = ResidueTypeStrand
			}
		}
	}
	return &Model{atoms, hetAtoms, connections, helixes, strands, residues, chains, nil, nil}
}

func (model *Model) RibbonMesh() *fauxgl.Mesh {
	mesh := fauxgl.NewEmptyMesh()
	for _, c := range model.Chains {
		m := c.Mesh()
		for i, t := range m.Triangles {
			p := float64(i) / float64(len(m.Triangles)-1)
			t.SetColor(fauxgl.MakeColor(Viridis.Color(p)))
		}
		mesh.Add(m)
	}
	return mesh
}

func (model *Model) HetMesh() *fauxgl.Mesh {
	mesh := fauxgl.NewEmptyMesh()
	sphere := fauxgl.NewSphere(15, 15)
	sphere.SmoothNormals()
	atomsBySerial := make(map[int]*Atom)
	for _, a := range model.HetAtoms {
		if a.ResName == "HOH" {
			continue
		}
		atomsBySerial[a.Serial] = a
		e := a.GetElement()
		r := e.Radius * 0.75
		s := fauxgl.V(r, r, r)
		m := sphere.Copy()
		m.Transform(fauxgl.Scale(s).Translate(a.Position))
		m.SetColor(fauxgl.HexColor(e.HexColor))
		mesh.Add(m)
	}
	for _, c := range model.Connections {
		a1 := atomsBySerial[c.Serial1]
		a2 := atomsBySerial[c.Serial2]
		mesh.Add(makeConnection(a1, a2))
	}
	return mesh
}

func (model *Model) SpaceFillingMesh() *fauxgl.Mesh {
	mesh := fauxgl.NewEmptyMesh()
	sphere := fauxgl.NewSphere(15, 15)
	sphere.SmoothNormals()
	for _, a := range model.Atoms {
		e := a.GetElement()
		r := e.VanDerWaalsRadius
		r = e.Radius * 0.5
		s := fauxgl.V(r, r, r)
		m := sphere.Copy()
		m.Transform(fauxgl.Scale(s).Translate(a.Position))
		m.SetColor(fauxgl.HexColor(e.HexColor))
		mesh.Add(m)
	}
	for _, r := range model.Residues {
		mesh.Add(makeConnection(r.Atoms["CA"], r.Atoms["N"]))
		mesh.Add(makeConnection(r.Atoms["C"], r.Atoms["O"]))
		mesh.Add(makeConnection(r.Atoms["C"], r.Atoms["CA"]))
		mesh.Add(makeConnection(r.Atoms["CA"], r.Atoms["CB"]))
		mesh.Add(makeConnection(r.Atoms["CB"], r.Atoms["CG"]))
		mesh.Add(makeConnection(r.Atoms["CG"], r.Atoms["CD"]))
		mesh.Add(makeConnection(r.Atoms["CD"], r.Atoms["CE"]))
	}
	return mesh
}

func (model *Model) Mesh() *fauxgl.Mesh {
	mesh := fauxgl.NewEmptyMesh()
	mesh.Add(model.RibbonMesh())
	mesh.Add(model.HetMesh())
	// mesh.Add(model.SpaceFillingMesh())

	// base := mesh.Copy()
	// for _, matrix := range model.SymmetryMatrixes {
	// 	if matrix == fauxgl.Identity() {
	// 		continue
	// 	}
	// 	m := base.Copy()
	// 	m.Transform(matrix)
	// 	mesh.Add(m)
	// }

	return mesh
}

func makeConnection(a1, a2 *Atom) *fauxgl.Mesh {
	mesh := fauxgl.NewEmptyMesh()
	if a1 == nil || a2 == nil {
		return mesh
	}
	e1 := a1.GetElement()
	e2 := a2.GetElement()
	p1 := a1.Position.LerpDistance(a2.Position, e1.Radius*0.75-0.1)
	p2 := a2.Position.LerpDistance(a1.Position, e2.Radius*0.75-0.1)
	mid := p1.Lerp(p2, 0.5)
	m := makeCylinder(p1, mid, 0.25)
	m.SetColor(fauxgl.HexColor(e1.HexColor))
	mesh.Add(m)
	m = makeCylinder(mid, p2, 0.25)
	m.SetColor(fauxgl.HexColor(e2.HexColor))
	mesh.Add(m)
	return mesh
}

func makeCylinder(p0, p1 fauxgl.Vector, r float64) *fauxgl.Mesh {
	p := p0.Add(p1).MulScalar(0.5)
	h := p0.Distance(p1) * 2
	up := p1.Sub(p0).Normalize()
	mesh := fauxgl.NewCylinder(15, false)
	mesh.Transform(fauxgl.Orient(p, fauxgl.V(r, r, h), up, 0))
	return mesh
}

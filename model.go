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
			if r.ChainID == h.ChainID && r.ResSeq >= h.InitSeqNum && r.ResSeq <= h.EndSeqNum {
				r.Type = ResidueTypeHelix
			}
		}
		for _, s := range strands {
			if r.ChainID == s.ChainID && r.ResSeq >= s.InitSeqNum && r.ResSeq <= s.EndSeqNum {
				r.Type = ResidueTypeStrand
			}
		}
	}
	return &Model{atoms, hetAtoms, connections, helixes, strands, residues, chains, nil, nil}
}

func (model *Model) Mesh() *fauxgl.Mesh {
	mesh := fauxgl.NewEmptyMesh()
	for _, c := range model.Chains {
		m := c.Mesh()
		for i, t := range m.Triangles {
			p := float64(i) / float64(len(m.Triangles)-1)
			t.SetColor(fauxgl.MakeColor(Viridis.Color(p)))
		}
		mesh.Add(m)
	}

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
		if a1 == nil || a2 == nil {
			continue
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
	}

	// base := mesh.Copy()
	// for _, matrix := range model.SymmetryMatrixes {
	// 	if matrix == Identity() {
	// 		continue
	// 	}
	// 	m := base.Copy()
	// 	m.Transform(matrix)
	// 	mesh.Add(m)
	// }

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

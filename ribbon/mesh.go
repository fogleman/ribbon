package ribbon

import (
	"github.com/fogleman/fauxgl"
	"github.com/fogleman/ribbon/pdb"
)

func ModelMesh(model *pdb.Model) *fauxgl.Mesh {
	mesh := fauxgl.NewEmptyMesh()
	mesh.Add(RibbonMesh(model))
	mesh.Add(HetMesh(model))
	// mesh.Add(SpaceFillingMesh(model))
	// mesh.Add(BackboneMesh(model))

	// base := mesh.Copy()
	// for _, m := range model.SymMatrixes {
	// 	mat := fauxgl.Matrix{
	// 		m[0][0], m[0][1], m[0][2], m[0][3],
	// 		m[1][0], m[1][1], m[1][2], m[1][3],
	// 		m[2][0], m[2][1], m[2][2], m[2][3],
	// 		m[3][0], m[3][1], m[3][2], m[3][3],
	// 	}
	// 	if mat == fauxgl.Identity() {
	// 		continue
	// 	}
	// 	meshCopy := base.Copy()
	// 	meshCopy.Transform(mat)
	// 	mesh.Add(meshCopy)
	// }

	return mesh
}

func RibbonMesh(model *pdb.Model) *fauxgl.Mesh {
	mesh := fauxgl.NewEmptyMesh()
	for _, chain := range model.Chains {
		m := createChainMesh(chain)
		for i, t := range m.Triangles {
			p := float64(i) / float64(len(m.Triangles)-1)
			t.SetColor(Viridis.Color(p))
		}
		mesh.Add(m)
	}
	// mesh.SmoothNormalsThreshold(fauxgl.Radians(60))
	return mesh
}

func HetMesh(model *pdb.Model) *fauxgl.Mesh {
	mesh := fauxgl.NewEmptyMesh()
	atomsBySerial := make(map[int]*pdb.Atom)
	for _, a := range model.HetAtoms {
		if a.ResName == "HOH" {
			continue
		}
		atomsBySerial[a.Serial] = a
		e := atomElement(a)
		r := e.Radius * 0.75
		s := fauxgl.V(r, r, r)
		m := unitSphere.Copy()
		m.Transform(fauxgl.Scale(s).Translate(atomPosition(a)))
		m.SetColor(fauxgl.HexColor(e.HexColor))
		mesh.Add(m)
	}
	for _, c := range model.Connections {
		if c.Serial1 > c.Serial2 {
			continue
		}
		a1 := atomsBySerial[c.Serial1]
		a2 := atomsBySerial[c.Serial2]
		mesh.Add(makeConnection(a1, a2))
	}
	return mesh
}

func SpaceFillingMesh(model *pdb.Model) *fauxgl.Mesh {
	mesh := fauxgl.NewEmptyMesh()
	for _, a := range model.Atoms {
		e := atomElement(a)
		r := e.VanDerWaalsRadius
		// r = e.Radius * 0.5
		s := fauxgl.V(r, r, r)
		m := unitSphere.Copy()
		m.Transform(fauxgl.Scale(s).Translate(atomPosition(a)))
		m.SetColor(fauxgl.HexColor(e.HexColor))
		mesh.Add(m)
	}
	// for _, r := range model.Residues {
	//  mesh.Add(makeConnection(r.Atoms["CA"], r.Atoms["N"]))
	//  mesh.Add(makeConnection(r.Atoms["C"], r.Atoms["O"]))
	//  mesh.Add(makeConnection(r.Atoms["C"], r.Atoms["CA"]))
	//  mesh.Add(makeConnection(r.Atoms["CA"], r.Atoms["CB"]))
	//  mesh.Add(makeConnection(r.Atoms["CB"], r.Atoms["CG"]))
	//  mesh.Add(makeConnection(r.Atoms["CG"], r.Atoms["CD"]))
	//  mesh.Add(makeConnection(r.Atoms["CD"], r.Atoms["CE"]))
	// }
	return mesh
}

func BackboneMesh(model *pdb.Model) *fauxgl.Mesh {
	mesh := fauxgl.NewEmptyMesh()
	for _, chain := range model.Chains {
		var previous *pdb.Residue
		for _, residue := range chain.Residues {
			var m *fauxgl.Mesh
			c1 := atomPosition(residue.AtomsByName["C"])
			m = makeSphere(c1, 0.4)
			m.SetColor(fauxgl.White)
			mesh.Add(m)
			o1 := atomPosition(residue.AtomsByName["O"])
			m = makeCylinder(c1, o1, 0.25)
			m.SetColor(fauxgl.HexColor(ElementsBySymbol["O"].HexColor))
			mesh.Add(m)
			m = makeSphere(o1, 0.25)
			m.SetColor(fauxgl.HexColor(ElementsBySymbol["O"].HexColor))
			mesh.Add(m)
			if previous != nil {
				c2 := atomPosition(previous.AtomsByName["C"])
				m = makeCylinder(c1, c2, 0.4)
				m.SetColor(fauxgl.White)
				mesh.Add(m)
			}
			previous = residue
		}
	}
	return mesh
}

func LineMesh(model *pdb.Model) *fauxgl.Mesh {
	var lines []*fauxgl.Line

	const n = 2
	for x := -n; x <= n; x++ {
		u := float64(x) / 4
		for _, chain := range model.Chains {
			var planes []*PeptidePlane
			for i := 0; i < len(chain.Residues)-2; i++ {
				r1 := chain.Residues[i]
				r2 := chain.Residues[i+1]
				r3 := chain.Residues[i+2]
				plane := NewPeptidePlane(r1, r2, r3)
				if plane != nil {
					planes = append(planes, plane)
				}
			}
			var previous fauxgl.Vector
			for i, p := range planes {
				if i > 0 && p.Side.Dot(previous) < 0 {
					p.Flip()
				}
				previous = p.Side
			}
			n := len(planes) - 3
			for i := 0; i < n; i++ {
				p1 := planes[i]
				p2 := planes[i+1]
				p3 := planes[i+2]
				p4 := planes[i+3]
				p := splineForPlanes(p1, p2, p3, p4, 4, u, 0)
				for j := 1; j < len(p); j++ {
					v1 := fauxgl.Vector{p[j-1].X, p[j-1].Y, p[j-1].Z}
					v2 := fauxgl.Vector{p[j].X, p[j].Y, p[j].Z}
					line := fauxgl.NewLineForPoints(v1, v2)
					lines = append(lines, line)
				}
			}
		}
	}

	return fauxgl.NewLineMesh(lines)
}

func makeConnection(a1, a2 *pdb.Atom) *fauxgl.Mesh {
	mesh := fauxgl.NewEmptyMesh()
	if a1 == nil || a2 == nil {
		return mesh
	}
	e1 := atomElement(a1)
	e2 := atomElement(a2)
	ap1 := atomPosition(a1)
	ap2 := atomPosition(a2)
	p1 := ap1.LerpDistance(ap2, e1.Radius*0.75-0.1)
	p2 := ap2.LerpDistance(ap1, e2.Radius*0.75-0.1)
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
	h := p0.Distance(p1)
	up := p1.Sub(p0).Normalize()
	mesh := unitCylinder.Copy()
	mesh.Transform(fauxgl.Orient(p, fauxgl.V(r, r, h), up, 0))
	return mesh
}

func makeSphere(p fauxgl.Vector, r float64) *fauxgl.Mesh {
	mesh := unitSphere.Copy()
	mesh.Transform(fauxgl.Scale(fauxgl.V(r, r, r)).Translate(p))
	return mesh
}

var (
	unitSphere   *fauxgl.Mesh
	unitCylinder *fauxgl.Mesh
)

func init() {
	unitSphere = fauxgl.NewSphere(15, 15)
	unitSphere.SmoothNormals()
	unitCylinder = fauxgl.NewCylinder(15, false)
	// unitCylinder.SmoothNormals()
}

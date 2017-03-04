package main

import (
	"fmt"
	"log"
	"os"
	"time"

	. "github.com/fogleman/fauxgl"
	"github.com/fogleman/ribbon"
	"github.com/nfnt/resize"
)

const (
	scale  = 4
	width  = 1600 * 1
	height = 1600 * 1
	fovy   = 30
	near   = 1
	far    = 10
)

// var (
// 	eye    = V(4, 0, 4)
// 	center = V(0, -0.03, 0)
// 	up     = V(0, 1, 0).Normalize()
// 	light  = V(0.75, 0.25, 0.25).Normalize()
// )

var (
	eye    = V(5, 0, 0)
	center = V(0, 0, 0)
	up     = V(0, 0, 1).Normalize()
	light  = V(0.75, 0.25, 0.25).Normalize()
)

func makeCylinder(p0, p1 Vector, r float64) *Mesh {
	p := p0.Add(p1).MulScalar(0.5)
	h := p0.Distance(p1) * 2
	up := p1.Sub(p0).Normalize()
	mesh := NewCylinder(15, false)
	mesh.Transform(Orient(p, V(r, r, h), up, 0))
	return mesh
}

func main() {
	model, err := ribbon.LoadPDB(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(model.Atoms), len(model.Residues), len(model.Chains))

	mesh := NewEmptyMesh()
	for _, c := range model.Chains {
		m := c.Mesh()
		mesh.Add(m)
	}
	fmt.Println(len(mesh.Triangles))

	sphere := NewSphere(15, 15)
	sphere.SmoothNormals()
	atomsBySerial := make(map[int]*ribbon.Atom)
	for _, a := range model.HetAtoms {
		if a.ResName == "HOH" {
			continue
		}
		atomsBySerial[a.Serial] = a
		e := a.GetElement()
		c := HexColor(e.HexColor)
		r := e.Radius * 0.75
		s := V(r, r, r)
		m := sphere.Copy()
		m.Transform(Scale(s).Translate(a.Position))
		for _, t := range m.Triangles {
			t.V1.Color = c
			t.V2.Color = c
			t.V3.Color = c
		}
		mesh.Add(m)
	}
	fmt.Println(len(mesh.Triangles))

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
		c := HexColor(e1.HexColor)
		for _, t := range m.Triangles {
			t.V1.Color = c
			t.V2.Color = c
			t.V3.Color = c
		}
		mesh.Add(m)
		m = makeCylinder(p2, mid, 0.25)
		c = HexColor(e2.HexColor)
		for _, t := range m.Triangles {
			t.V1.Color = c
			t.V2.Color = c
			t.V3.Color = c
		}
		mesh.Add(m)
	}
	fmt.Println(len(mesh.Triangles))

	// base := mesh.Copy()
	// for _, matrix := range model.SymmetryMatrixes {
	// 	if matrix == Identity() {
	// 		continue
	// 	}
	// 	m := base.Copy()
	// 	m.Transform(matrix)
	// 	mesh.Add(m)
	// }
	// fmt.Println(len(mesh.Triangles))

	mesh.BiUnitCube()
	// mesh.SmoothNormalsThreshold(Radians(75))
	// mesh.SaveSTL("out.stl")

	// create a rendering context
	context := NewContext(width*scale, height*scale)
	context.ClearColorBufferWith(HexColor("1D181F"))

	// create transformation matrix and light direction
	aspect := float64(width) / float64(height)
	matrix := LookAt(eye, center, up).Perspective(fovy, aspect, near, far)

	// render
	shader := NewPhongShader(matrix, light, eye)
	shader.AmbientColor = Gray(0.3)
	shader.DiffuseColor = Gray(0.9)
	context.Shader = shader
	// context.Cull = CullFront
	start := time.Now()
	context.DrawTriangles(mesh.Triangles)
	fmt.Println(time.Since(start))

	// context.ClearDepthBuffer()
	// start = time.Now()
	// context.Cull = CullBack
	// context.DrawTriangles(mesh.Triangles)
	// fmt.Println(time.Since(start))

	// save image
	image := context.Image()
	image = resize.Resize(width, height, image, resize.Bilinear)
	SavePNG("out.png", image)

	// for i := 0; i < 360; i += 1 {
	// 	context.ClearColorBufferWith(HexColor("1D181F"))
	// 	context.ClearDepthBuffer()

	// 	shader := NewPhongShader(matrix, light, eye)
	// 	shader.AmbientColor = Gray(0.3)
	// 	shader.DiffuseColor = Gray(0.9)
	// 	context.Shader = shader
	// 	start := time.Now()
	// 	context.DepthBias = 0
	// 	context.DrawTriangles(mesh.Triangles)
	// 	fmt.Println(time.Since(start))

	// 	image := context.Image()
	// 	image = resize.Resize(width, height, image, resize.Bilinear)
	// 	SavePNG(fmt.Sprintf("frame%03d.png", i), image)

	// 	mesh.Transform(Rotate(up, Radians(1)))
	// }
}

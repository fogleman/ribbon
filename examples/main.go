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
	width  = 1920
	height = 1080
	fovy   = 25
	near   = 1
	far    = 10
)

// var (
// 	eye    = V(0, -4, 0)
// 	center = V(0.15, 0, 0.15)
// 	up     = V(-1, 0, -1).Normalize()
// 	light  = V(0.25, -0.75, 0.25).Normalize()
// )

var (
	eye    = V(4, 0, 0)
	center = V(0, 0, 0)
	up     = V(0, 1, 1).Normalize()
	light  = V(0.75, 0.25, 0.25).Normalize()
)

func makeCylinder(p0, p1 Vector, r float64) *Mesh {
	p := p0.Add(p1).MulScalar(0.5)
	h := p0.Distance(p1) * 2
	up := p1.Sub(p0).Normalize()
	mesh := NewCylinder(30, false)
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
		// color := Color{rand.Float64(), rand.Float64(), rand.Float64(), 1}
		// for _, t := range m.Triangles {
		// 	t.V1.Color = color
		// 	t.V2.Color = color
		// 	t.V3.Color = color
		// }
		mesh.Add(m)
		// break
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
		e := ribbon.ElementsBySymbol[a.Element]
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
		a := atomsBySerial[c.Serial1]
		b := atomsBySerial[c.Serial2]
		if a == nil || b == nil {
			continue
		}
		m := makeCylinder(a.Position, b.Position, 0.25)
		for _, t := range m.Triangles {
			t.V1.Color = White
			t.V2.Color = White
			t.V3.Color = White
		}
		mesh.Add(m)
	}
	fmt.Println(len(mesh.Triangles))

	// var previous Vector
	// for i, r := range model.Residues {
	// 	a := r.Atoms["CA"]
	// 	if a == nil {
	// 		continue
	// 	}
	// 	if a.ChainID != "A" {
	// 		continue
	// 	}
	// 	color := White
	// 	if r.Type == ribbon.ResidueTypeHelix {
	// 		color = Color{1, 0, 0, 1}
	// 	}
	// 	if r.Type == ribbon.ResidueTypeStrand {
	// 		color = Color{0, 1, 0, 1}
	// 	}
	// 	s := NewSphere(15, 15)
	// 	for _, t := range s.Triangles {
	// 		t.V1.Color = color
	// 		t.V2.Color = color
	// 		t.V3.Color = color
	// 	}
	// 	s.Transform(Scale(V(0.333, 0.333, 0.333)).Translate(a.Position))
	// 	mesh.Add(s)
	// 	if i != 0 {
	// 		c := makeCylinder(previous, a.Position, 0.2)
	// 		for _, t := range c.Triangles {
	// 			t.V1.Color = White
	// 			t.V2.Color = White
	// 			t.V3.Color = White
	// 		}
	// 		mesh.Add(c)
	// 	}
	// 	previous = a.Position
	// }

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
	context.Cull = CullFront
	start := time.Now()
	context.DrawTriangles(mesh.Triangles)
	fmt.Println(time.Since(start))

	context.ClearDepthBuffer()
	start = time.Now()
	context.Cull = CullBack
	context.DrawTriangles(mesh.Triangles)
	fmt.Println(time.Since(start))

	// save image
	image := context.Image()
	image = resize.Resize(width, height, image, resize.Bilinear)
	SavePNG("out.png", image)

	// for i := 0; i < 360; i += 1 {
	// 	context.ClearColorBufferWith(HexColor("2A2C2B"))
	// 	context.ClearDepthBuffer()

	// 	shader := NewPhongShader(matrix, light, eye)
	// 	shader.AmbientColor = Gray(0.3)
	// 	shader.DiffuseColor = Gray(0.9)
	// 	context.Shader = shader
	// 	context.Cull = CullNone
	// 	start := time.Now()
	// 	context.DepthBias = 0
	// 	context.DrawTriangles(mesh.Triangles)
	// 	fmt.Println(time.Since(start))

	// 	context.Shader = NewSolidColorShader(matrix, Black)
	// 	context.LineWidth = scale * 1.5
	// 	context.DepthBias = -1e-4
	// 	context.DrawLines(mesh.Lines)

	// 	image := context.Image()
	// 	image = resize.Resize(width, height, image, resize.Bilinear)
	// 	SavePNG(fmt.Sprintf("frame%03d.png", i), image)

	// 	mesh.Transform(Rotate(up, Radians(1)))
	// }
}

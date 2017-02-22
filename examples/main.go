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
	width  = 2048
	height = 2048
	fovy   = 50
	near   = 1
	far    = 10
)

var (
	eye    = V(3, 0, 0)
	center = V(0, 0, 0)
	up     = V(0, 0, 1)
	light  = V(0.75, 0.25, 0.25).Normalize()
)

var colors = []Color{
	HexColor("7F1637"),
	HexColor("047878"),
	HexColor("FFB733"),
	HexColor("F57336"),
	HexColor("C22121"),
}

func main() {
	model, err := ribbon.LoadPDB(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(model.Atoms), len(model.Residues), len(model.Chains))

	mesh := NewEmptyMesh()
	for i, c := range model.Chains {
		m := c.Ribbon(3, 0.1)
		c := colors[i%len(colors)]
		for _, t := range m.Triangles {
			t.V1.Color = c
			t.V2.Color = c
			t.V3.Color = c
		}
		mesh.Add(m)
	}
	mesh.BiUnitCube()
	// mesh.SmoothNormalsThreshold(Radians(90))

	// var edges []*Triangle
	// for _, t := range mesh.Triangles {
	// 	n := t.Normal()
	// 	e1 := math.Abs(t.V1.Position.Sub(eye).Normalize().Dot(n)) < 0.05
	// 	e2 := math.Abs(t.V2.Position.Sub(eye).Normalize().Dot(n)) < 0.05
	// 	e3 := math.Abs(t.V3.Position.Sub(eye).Normalize().Dot(n)) < 0.05
	// 	if e1 && e2 && e3 {
	// 		edges = append(edges, t)
	// 	}
	// }

	// create a rendering context
	context := NewContext(width*scale, height*scale)
	context.ClearColorBufferWith(HexColor("323"))

	// create transformation matrix and light direction
	aspect := float64(width) / float64(height)
	matrix := LookAt(eye, center, up).Perspective(fovy, aspect, near, far)

	// render
	shader := NewPhongShader(matrix, light, eye)
	shader.AmbientColor = Gray(0.3)
	shader.DiffuseColor = Gray(0.9)
	context.Shader = shader
	start := time.Now()
	context.DrawTriangles(mesh.Triangles)
	fmt.Println(time.Since(start))

	context.Shader = NewSolidColorShader(matrix, Black)
	context.LineWidth = scale * 1.5
	context.DepthBias = -5e-4
	context.DrawLines(mesh.Lines)

	// save image
	image := context.Image()
	image = resize.Resize(width, height, image, resize.Bilinear)
	SavePNG("out.png", image)

	// for i := 0; i < 360; i += 1 {
	// 	context.ClearColorBufferWith(HexColor("323"))
	// 	context.ClearDepthBuffer()

	// 	start := time.Now()
	// 	context.DrawMesh(mesh)
	// 	fmt.Println(time.Since(start))

	// 	image := context.Image()
	// 	image = resize.Resize(width, height, image, resize.Bilinear)
	// 	SavePNG(fmt.Sprintf("frame%03d.png", i), image)

	// 	mesh.Transform(Rotate(up, Radians(1)))
	// }
}

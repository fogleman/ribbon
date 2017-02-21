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
	width  = 1600 * 2
	height = 1200 * 2
	fovy   = 35
	near   = 1
	far    = 10
)

var (
	eye    = V(3, 1, 0.5)
	center = V(0, -0.1, 0)
	up     = V(0, 0, 1)
	light  = V(0.75, 0.25, 1).Normalize()
)

func main() {
	model, err := ribbon.LoadPDB(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(model.Atoms), len(model.Residues), len(model.Polypeptides))

	mesh := NewEmptyMesh()
	for _, pp := range model.Polypeptides {
		mesh.Add(pp.Ribbon(0, 0))
	}
	mesh.BiUnitCube()
	mesh.SmoothNormalsThreshold(Radians(60))
	fmt.Println(len(mesh.Triangles))

	// create a rendering context
	context := NewContext(width*scale, height*scale)
	context.ClearColorBufferWith(HexColor("323"))

	// create transformation matrix and light direction
	aspect := float64(width) / float64(height)
	matrix := LookAt(eye, center, up).Perspective(fovy, aspect, near, far)

	// render
	shader := NewPhongShader(matrix, light, eye)
	// shader.ObjectColor = HexColor("FFD34E")
	context.Shader = shader
	start := time.Now()
	// context.Cull = CullNone
	context.DrawMesh(mesh)
	fmt.Println(time.Since(start))

	// save image
	image := context.Image()
	image = resize.Resize(width, height, image, resize.Bilinear)
	SavePNG("out.png", image)
}

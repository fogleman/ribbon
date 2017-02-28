package main

import (
	"log"

	. "github.com/fogleman/fauxgl"
	"github.com/fogleman/ribbon"
	"github.com/nfnt/resize"
)

const (
	scale  = 1
	width  = 4000 * 1
	height = 1000 * 1
	fovy   = 8
	near   = 1
	far    = 10
)

var (
	eye    = V(0, 2, 4)
	center = V(0, 0, 0)
	up     = V(0, 1, 0).Normalize()
	light  = V(0.25, 0.25, 0.75).Normalize()
)

func main() {
	model, err := ribbon.LoadPDB("examples/test.pdb")
	if err != nil {
		log.Fatal(err)
	}

	mesh := NewEmptyMesh()
	for _, c := range model.Chains {
		m := c.Mesh()
		mesh.Add(m)
	}

	mesh.BiUnitCube()
	// mesh.SmoothNormalsThreshold(Radians(75))
	mesh.SaveSTL("out.stl")

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
	context.DrawTriangles(mesh.Triangles)

	// save image
	image := context.Image()
	image = resize.Resize(width, height, image, resize.Bilinear)
	SavePNG("out.png", image)
}

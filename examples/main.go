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
	width  = 2400 * 1
	height = 1600 * 1
	fovy   = 25
	near   = 1
	far    = 10
)

var (
	eye    = V(5, 0, 0)
	center = V(0, 0, 0)
	up     = V(0, 1, 0).Normalize()
	light  = eye.Sub(center).Normalize()
)

func dumpMesh(mesh *Mesh) {
	var vertices []Vector
	var colors []Color
	indexLookup := make(map[Vector]int)
	colorLookup := make(map[Color]int)
	for _, t := range mesh.Triangles {
		if _, ok := indexLookup[t.V1.Position]; !ok {
			indexLookup[t.V1.Position] = len(vertices)
			vertices = append(vertices, t.V1.Position)
		}
		if _, ok := indexLookup[t.V2.Position]; !ok {
			indexLookup[t.V2.Position] = len(vertices)
			vertices = append(vertices, t.V2.Position)
		}
		if _, ok := indexLookup[t.V3.Position]; !ok {
			indexLookup[t.V3.Position] = len(vertices)
			vertices = append(vertices, t.V3.Position)
		}
		if _, ok := colorLookup[t.V1.Color]; !ok {
			colorLookup[t.V1.Color] = len(colors)
			colors = append(colors, t.V1.Color)
		}
		if _, ok := colorLookup[t.V2.Color]; !ok {
			colorLookup[t.V2.Color] = len(colors)
			colors = append(colors, t.V2.Color)
		}
		if _, ok := colorLookup[t.V3.Color]; !ok {
			colorLookup[t.V3.Color] = len(colors)
			colors = append(colors, t.V3.Color)
		}
	}
	fmt.Println("var VERTICES = [")
	for _, v := range vertices {
		fmt.Printf("[%.8f,%.8f,%.8f],\n", v.X, v.Y, v.Z)
	}
	fmt.Println("];")
	fmt.Println("var COLORS = [")
	for _, c := range colors {
		fmt.Printf("[%.3f,%.3f,%.3f],\n", c.R, c.G, c.B)
	}
	fmt.Println("];")
	fmt.Println("var FACES = [")
	for _, t := range mesh.Triangles {
		i1 := indexLookup[t.V1.Position]
		i2 := indexLookup[t.V2.Position]
		i3 := indexLookup[t.V3.Position]
		c1 := colorLookup[t.V1.Color]
		c2 := colorLookup[t.V2.Color]
		c3 := colorLookup[t.V3.Color]
		fmt.Printf("[%d,%d,%d,%d,%d,%d],\n", i1, i2, i3, c1, c2, c3)
	}
	fmt.Println("];")
}

func main() {
	model, err := ribbon.LoadPDB(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(len(model.Atoms), len(model.Residues), len(model.Chains))

	mesh := model.Mesh()
	mesh.BiUnitCube()
	// dumpMesh(mesh)
	// return

	// mesh.SaveSTL("out.stl")

	// render
	context := NewContext(width*scale, height*scale)
	context.ClearColorBufferWith(HexColor("1D181F"))
	aspect := float64(width) / float64(height)
	matrix := LookAt(eye, center, up).Perspective(fovy, aspect, near, far)
	shader := NewPhongShader(matrix, light, eye)
	shader.AmbientColor = Gray(0.3)
	shader.DiffuseColor = Gray(0.9)
	context.Shader = shader
	start := time.Now()
	context.DrawTriangles(mesh.Triangles)
	fmt.Println(time.Since(start))

	// save image
	image := context.Image()
	image = resize.Resize(width, height, image, resize.Bilinear)
	SavePNG("out.png", image)

	// for i := 0; i < 720; i += 1 {
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

	// 	mesh.Transform(Rotate(up, Radians(0.5)))
	// }
}

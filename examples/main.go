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
	width  = 2048 * 1
	height = 2048 * 1
	near   = 1
	far    = 100
)

var (
	eye    = V(5, 0, 0)
	center = V(0, 0, 0)
	up     = V(0, 1, 0).Normalize()
	fovy   = 20.0
	light  = eye.Sub(center).Normalize()
)

func timed(name string) func() {
	if len(name) > 0 {
		fmt.Printf("%s... ", name)
	}
	start := time.Now()
	return func() {
		fmt.Println(time.Since(start))
	}
}

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
	var done func()

	done = timed("loading pdb file")
	model, err := ribbon.LoadPDB(os.Args[1])
	done()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("atoms       = %d\n", len(model.Atoms))
	fmt.Printf("residues    = %d\n", len(model.Residues))
	fmt.Printf("chains      = %d\n", len(model.Chains))
	fmt.Printf("helixes     = %d\n", len(model.Helixes))
	fmt.Printf("strands     = %d\n", len(model.Strands))
	fmt.Printf("het-atoms   = %d\n", len(model.HetAtoms))
	fmt.Printf("connections = %d\n", len(model.Connections))

	// min := model.Atoms[0].TempFactor
	// max := model.Atoms[0].TempFactor
	// for _, a := range model.Atoms {
	// 	if a.Name != "CA" {
	// 		continue
	// 	}
	// 	min = math.Min(min, a.TempFactor)
	// 	max = math.Max(max, a.TempFactor)
	// }
	// fmt.Println(min, max)

	done = timed("generating triangle mesh")
	mesh := model.Mesh()
	done()

	fmt.Printf("triangles   = %d\n", len(mesh.Triangles))

	done = timed("transforming mesh")
	m := mesh.BiUnitCube()
	done()
	// dumpMesh(mesh)
	// return

	done = timed("finding ideal camera position")
	camera := model.Camera(m)
	eye = camera.Eye
	center = camera.Center
	up = camera.Up
	fovy = camera.Fovy
	light = eye.Sub(center).Normalize()
	done()

	// mesh.SaveSTL("out.stl")

	// render
	done = timed("rendering image")
	context := NewContext(width*scale, height*scale)
	context.ClearColorBufferWith(HexColor("1D181F"))
	aspect := float64(width) / float64(height)
	matrix := LookAt(eye, center, up).Perspective(fovy, aspect, near, far)
	shader := NewPhongShader(matrix, light, eye)
	shader.AmbientColor = Gray(0.3)
	shader.DiffuseColor = Gray(0.9)
	context.Shader = shader
	context.DrawTriangles(mesh.Triangles)
	done()

	// save image
	done = timed("downsampling image")
	image := context.Image()
	image = resize.Resize(width, height, image, resize.Bilinear)
	done()

	done = timed("writing image to disk")
	// SavePNG(os.Args[2], image)
	SavePNG("out.png", image)
	done()

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

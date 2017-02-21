package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	. "github.com/fogleman/fauxgl"
	"github.com/fogleman/ribbon"
	"github.com/nfnt/resize"
)

const (
	scale  = 8
	width  = 1600
	height = 1200
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

func fv(v ribbon.Vector) Vector {
	return Vector{v.X, v.Y, v.Z}
}

func main() {
	model, err := ribbon.LoadPDB(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(model.Atoms), len(model.Residues))

	mesh := NewEmptyMesh()
	var previous *ribbon.ResiduePlane
	for i := 0; i < len(model.Residues)-1; i++ {
		r1 := model.Residues[i]
		r2 := model.Residues[i+1]
		if r1.Atoms["CA"] == nil || r2.Atoms["CA"] == nil {
			continue
		}
		if r1.Atoms["CA"].Position.Distance(r2.Atoms["CA"].Position) > 5 {
			continue
		}
		p := ribbon.NewResiduePlane(r1, r2)
		ps := []*ribbon.ResiduePlane{p}
		if previous != nil {
			if p.Position.Distance(previous.Position) < 5 {
				for i := 1; i < 32; i++ {
					ps = append(ps, previous.Lerp(p, float64(i)/32))
				}
			}
		}
		for _, p := range ps {
			plane := NewPlane()
			r := RotateTo(fv(p.Normal), up).MulPosition(fv(p.Forward))
			plane.Transform(Orient(fv(p.Position), V(2, 2, 2), fv(p.Normal), math.Atan2(r.Y, r.X)))
			mesh.Add(plane)
		}
		previous = p
	}

	mesh.BiUnitCube()

	// create a rendering context
	context := NewContext(width*scale, height*scale)
	context.ClearColorBufferWith(HexColor("323"))

	// create transformation matrix and light direction
	aspect := float64(width) / float64(height)
	matrix := LookAt(eye, center, up).Perspective(fovy, aspect, near, far)

	// render
	shader := NewPhongShader(matrix, light, eye)
	shader.ObjectColor = HexColor("FFD34E")
	context.Shader = shader
	start := time.Now()
	context.Cull = CullNone
	context.DrawMesh(mesh)
	fmt.Println(time.Since(start))

	// save image
	image := context.Image()
	image = resize.Resize(width, height, image, resize.Bilinear)
	SavePNG("out.png", image)
}

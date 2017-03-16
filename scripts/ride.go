package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"time"

	. "github.com/fogleman/fauxgl"
	"github.com/fogleman/ribbon"
)

const (
	scale  = 1
	width  = 1920
	height = 1080
	fovy   = 60
	near   = 0.1
	far    = 100
)

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
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

	ps, us := model.Chains[0].Poses(16)

	context := NewContext(width*scale, height*scale)

	for i := 0; i < len(ps)-16; i++ {
		path := fmt.Sprintf("frames/frame%06d.png", i)
		if exists(path) {
			continue
		}
		fmt.Println(i, len(ps))

		eye1 := ps[i]
		eye2 := ps[i+16]
		up1 := us[i].Sub(eye1).Normalize()
		up2 := us[i+16].Sub(eye2).Normalize()

		eye1 = eye1.Add(up1.MulScalar(0.3))
		eye2 = eye2.Add(up2.MulScalar(0.3))

		eye := eye1
		center := eye2
		up := up1
		light := eye1.Sub(eye2).Normalize()

		context.ClearColorBufferWith(HexColor("1D181F"))
		context.ClearDepthBuffer()
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

		// save image
		im := context.Image().(*image.NRGBA)
		im2 := image.NewNRGBA(im.Bounds())
		copy(im2.Pix, im.Pix)
		// im = resize.Resize(width, height, im, resize.Bilinear)
		go func(i int, im *image.NRGBA) {
			SavePNG(path, im)
		}(i, im2)
	}
}

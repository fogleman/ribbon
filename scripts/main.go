package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	. "github.com/fogleman/fauxgl"
	"github.com/fogleman/ribbon/pdb"
	"github.com/fogleman/ribbon/ribbon"
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
	fovy   = 30.0
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

func downloadModels(structureID string) ([]*pdb.Model, error) {
	url := fmt.Sprintf(
		"https://files.rcsb.org/download/%s.pdb.gz",
		strings.ToUpper(structureID))
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	gr, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, err
	}
	r := pdb.NewReader(gr)
	return r.ReadAll()
}

func loadStructureIDs(path string) ([]string, error) {
	var result []string
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func run(structureID string) {
	fmt.Println(structureID)

	var done func()

	// done = timed("loading pdb file")
	// file, err := os.Open(os.Args[1])
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// r := pdb.NewReader(file)
	// model, err := r.Read()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// done()

	done = timed("downloading pdb file")
	models, err := downloadModels(structureID)
	if err != nil {
		log.Fatal(err)
	}
	model := models[0]
	done()

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
	mesh := ribbon.ModelMesh(model)
	done()

	fmt.Printf("triangles   = %d\n", len(mesh.Triangles))

	done = timed("transforming mesh")
	m := mesh.BiUnitCube()
	done()
	// dumpMesh(mesh)
	// return

	done = timed("finding ideal camera position")
	camera := ribbon.PositionCamera(model, m)
	eye = camera.Eye
	center = camera.Center
	up = camera.Up
	fovy = camera.Fovy
	light = eye.Sub(center).Normalize()
	aspect := camera.Aspect
	done()

	// done = timed("writing mesh to disk")
	// mesh.SaveSTL(fmt.Sprintf("stl/%s.stl", structureID))
	// done()

	// render
	done = timed("rendering image")
	context := NewContext(int(width*scale*aspect), height*scale)
	context.ClearColorBufferWith(HexColor("1D181F"))
	// aspect := float64(width) / float64(height)
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
	image = resize.Resize(uint(width*aspect), height, image, resize.Bilinear)
	done()

	done = timed("writing image to disk")
	// SavePNG("out.png", image)
	SavePNG(fmt.Sprintf("png/%s.png", structureID), image)
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

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	structures, err := loadStructureIDs("structures.txt")
	if err != nil {
		log.Fatal(err)
	}
	for _, i := range rand.Perm(len(structures)) {
		run(structures[i])
	}
}

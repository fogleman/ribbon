package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	. "github.com/fogleman/fauxgl"
	"github.com/fogleman/mc"
	"github.com/fogleman/ribbon/pdb"
	"github.com/fogleman/ribbon/ribbon"
)

const (
	voxelSizeAngstroms = 1.
	sigmaAngstroms     = 4.
	thresholdAngstroms = 4.
	truncate           = 3.

	sigmaVoxels     = sigmaAngstroms / voxelSizeAngstroms
	thresholdVoxels = thresholdAngstroms / voxelSizeAngstroms
)

func newGaussianKernel3D(sigma, standardDeviations float64) (int, []float64) {
	r := int(math.Ceil(sigma * standardDeviations))
	w := 2*r + 1

	k := make([]float64, w)
	for i := range k {
		x := float64(i - r)
		k[i] = math.Exp(-x * x / (2 * sigma * sigma))
	}

	kernel := make([]float64, w*w*w)

	var sum float64
	for z, wz := range k {
		for y, wy := range k {
			for x, wx := range k {
				i := x + y*w + z*w*w
				w := wx * wy * wz
				kernel[i] = w
				sum += w
			}
		}
	}

	// for i, w := range kernel {
	// 	kernel[i] = w / sum
	// }

	return w, kernel
}

func main() {
	args := os.Args[1:]
	if len(args) != 1 || len(args[0]) != 4 {
		fmt.Println("Usage: surface XXXX")
		fmt.Println(" XXXX: 4-digit RCSB PDB Structure ID")
		os.Exit(1)
	}
	structureID := args[0]

	var done func()

	done = timed("downloading pdb file")
	models, err := downloadAndParse(structureID)
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
	fmt.Printf("biomatrixes = %d\n", len(model.BioMatrixes))
	fmt.Printf("symmatrixes = %d\n", len(model.SymMatrixes))

	// get atom positions and radii
	spheres := ribbon.Spheres(model)

	// done = timed("downloading pdb file")
	// spheres, err := downloadAndParseCIF(structureID)
	// done()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	done = timed("computing bounds")

	// get bounding box of spheres
	lo := spheres[0].Vector()
	hi := spheres[0].Vector()
	for _, s := range spheres {
		lo = lo.Min(s.Vector())
		hi = hi.Max(s.Vector())
	}
	size := hi.Sub(lo)

	// compute kernel
	sigma := sigmaVoxels
	if sigma <= 0 {
		sigma = math.Pow(size.X*size.Y*size.Z, 1.0/9) / voxelSizeAngstroms
	}
	kw, kernel := newGaussianKernel3D(sigma, truncate)
	kr := kw / 2

	// compute voxel bounds / sizes
	gw := int(math.Ceil(size.X/voxelSizeAngstroms)) + kw
	gh := int(math.Ceil(size.Y/voxelSizeAngstroms)) + kw
	gd := int(math.Ceil(size.Z/voxelSizeAngstroms)) + kw
	grid := make([]float64, gw*gh*gd)

	done()

	fmt.Println(size)

	// apply kernel
	done = timed("applying kernel")
	for _, s := range spheres {
		x0 := int((s.X - lo.X) / voxelSizeAngstroms)
		y0 := int((s.Y - lo.Y) / voxelSizeAngstroms)
		z0 := int((s.Z - lo.Z) / voxelSizeAngstroms)
		for dz := 0; dz < kw; dz++ {
			for dy := 0; dy < kw; dy++ {
				for dx := 0; dx < kw; dx++ {
					x := x0 + dx
					y := y0 + dy
					z := z0 + dz
					gi := x + y*gw + z*gw*gh
					ki := dx + dy*kw + dz*kw*kw
					grid[gi] += kernel[ki]
				}
			}
		}
	}
	done()

	// compute threshold
	var sum, count float64
	for _, s := range spheres {
		x := int((s.X-lo.X)/voxelSizeAngstroms) + kr
		y := int((s.Y-lo.Y)/voxelSizeAngstroms) + kr
		z := int((s.Z-lo.Z)/voxelSizeAngstroms) + kr
		sum += grid[x+y*gw+z*gw*gh]
		count++
	}
	pct := math.Exp(-thresholdVoxels * thresholdVoxels / (2 * sigma * sigma))
	threshold := sum / count * pct

	// run marching cubes
	done = timed("running marching cubes")
	mcTriangles := mc.MarchingCubesGrid(gw, gh, gd, grid, threshold)
	done()

	// convert to mesh
	done = timed("converting to mesh")
	triangles := make([]*Triangle, len(mcTriangles))
	transform := Translate(lo).Translate(V(-float64(kr), -float64(kr), -float64(kr))).Scale(V(voxelSizeAngstroms, voxelSizeAngstroms, voxelSizeAngstroms))
	for i, t := range mcTriangles {
		p1 := transform.MulPosition(Vector(t.V1))
		p2 := transform.MulPosition(Vector(t.V2))
		p3 := transform.MulPosition(Vector(t.V3))
		triangles[i] = NewTriangleForPoints(p1, p2, p3)
	}
	mesh := NewTriangleMesh(triangles)

	// sphere := NewSphere(1)
	// for _, s := range spheres {
	// 	m := sphere.Copy()
	// 	m.Transform(Translate(s.Vector()))
	// 	mesh.Add(m)
	// }

	done()

	done = timed("writing mesh to disk")
	mesh.SaveSTL(fmt.Sprintf("%s.stl", structureID))
	done()

	// for z := 0; z < gd; z++ {
	// 	im := image.NewGray(image.Rect(0, 0, gw, gh))
	// 	slice := grid[z*(gw*gh) : (z+1)*(gw*gh)]
	// 	for i, v := range slice {
	// 		im.Pix[i] = uint8(v / maxGridValue * 255)
	// 	}
	// 	gg.SavePNG(fmt.Sprintf("%08d.png", z), im)
	// }
}

func downloadAndParse(structureID string) ([]*pdb.Model, error) {
	url := fmt.Sprintf(
		"https://files.rcsb.org/download/%s.pdb.gz",
		strings.ToUpper(structureID))
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	r, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, err
	}
	return pdb.NewReader(r).ReadAll()
}

func downloadAndParseCIF(structureID string) ([]VectorW, error) {
	url := fmt.Sprintf(
		"https://files.rcsb.org/download/%s.cif.gz",
		strings.ToUpper(structureID))
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	r, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, err
	}
	var result []VectorW
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		// fmt.Println(scanner.Text()) // Println will add back the final '\n'
		fields := strings.Fields(scanner.Text())
		if fields[0] != "ATOM" {
			continue
		}
		x, _ := strconv.ParseFloat(fields[10], 64)
		y, _ := strconv.ParseFloat(fields[11], 64)
		z, _ := strconv.ParseFloat(fields[12], 64)
		result = append(result, VectorW{x, y, z, 0})
	}
	return result, scanner.Err()
}

func timed(name string) func() {
	if len(name) > 0 {
		fmt.Printf("%s... ", name)
	}
	start := time.Now()
	return func() {
		fmt.Println(time.Since(start))
	}
}

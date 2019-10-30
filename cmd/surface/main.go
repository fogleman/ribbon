package main

import (
	"compress/gzip"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	. "github.com/fogleman/fauxgl"
	"github.com/fogleman/mc"
	"github.com/fogleman/ribbon/pdb"
	"github.com/fogleman/ribbon/ribbon"
)

const (
	voxelSizeAngstroms = 0.5
	sigmaAngstroms     = 4.
	thresholdAngstroms = 2.
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

	for i, w := range kernel {
		kernel[i] = w / sum
	}

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

	done = timed("computing bounds")

	// get atom positions and radii
	spheres := ribbon.Spheres(model)

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
		x := int((s.X - lo.X) / voxelSizeAngstroms)
		y := int((s.Y - lo.Y) / voxelSizeAngstroms)
		z := int((s.Z - lo.Z) / voxelSizeAngstroms)
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
	for i, t := range mcTriangles {
		p1 := Vector(t.V1).MulScalar(voxelSizeAngstroms)
		p2 := Vector(t.V2).MulScalar(voxelSizeAngstroms)
		p3 := Vector(t.V3).MulScalar(voxelSizeAngstroms)
		triangles[i] = NewTriangleForPoints(p1, p2, p3)
	}
	mesh := NewTriangleMesh(triangles)
	// mesh.Transform(Rotate(V(1, 0, 0), Radians(-135)))
	done()

	done = timed("writing mesh to disk")
	mesh.SaveSTL(fmt.Sprintf("%s.stl", structureID))
	done()
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

func timed(name string) func() {
	if len(name) > 0 {
		fmt.Printf("%s... ", name)
	}
	start := time.Now()
	return func() {
		fmt.Println(time.Since(start))
	}
}

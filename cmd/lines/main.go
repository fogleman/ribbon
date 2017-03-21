package main

import (
	"compress/gzip"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	. "github.com/fogleman/fauxgl"
	"github.com/fogleman/ribbon/pdb"
	"github.com/fogleman/ribbon/ribbon"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 || len(args[0]) != 4 {
		fmt.Println("Usage: rcsb XXXX")
		fmt.Println(" XXXX: 4-digit RCSB PDB Structure ID")
		os.Exit(1)
	}
	structureID := args[0]

	models, err := downloadAndParse(structureID)
	if err != nil {
		log.Fatal(err)
	}
	model := models[0]
	mesh := ribbon.LineMesh(model)
	m := mesh.BiUnitCube()
	c := ribbon.PositionCamera(model, m)

	matrix := LookAt(c.Eye, c.Center, c.Up).Perspective(c.Fovy, c.Aspect, 1, 100)
	for _, line := range mesh.Lines {
		v1 := matrix.MulPositionW(line.V1.Position)
		v1 = v1.DivScalar(v1.W)
		v2 := matrix.MulPositionW(line.V2.Position)
		v2 = v2.DivScalar(v2.W)
		fmt.Printf("[(%g, %g), (%g, %g)],\n", v1.X, v1.Y, v2.X, v2.Y)
	}
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

package pdb

import (
	"fmt"
	"os"
	"testing"
)

func TestReader(t *testing.T) {
	file, err := os.Open("5ujw.pdb")
	if err != nil {
		t.Error(err)
	}
	r := NewReader(file)
	models, err := r.ReadAll()
	if err != nil {
		t.Error(err)
	}
	for _, model := range models {
		fmt.Println(len(model.Atoms))
		fmt.Println(len(model.HetAtoms))
		fmt.Println(len(model.Connections))
		fmt.Println(len(model.Helixes))
		fmt.Println(len(model.Strands))
		fmt.Println(len(model.BioMatrixes))
		fmt.Println(len(model.SymMatrixes))
		fmt.Println()
	}
}

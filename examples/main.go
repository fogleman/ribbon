package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fogleman/ribbon"
)

func main() {
	model, err := ribbon.LoadPDB(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(model.Atoms), len(model.Residues))
	for _, r := range model.Residues {
		fmt.Println(r.Name, r.Atoms["CA"], r.Atoms["O"])
	}
}

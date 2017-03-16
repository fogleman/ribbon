package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "net/http/pprof"

	"github.com/fogleman/ribbon/pdb"
	"github.com/fogleman/ribbon/ribbon"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	var start time.Time
	for {
		file, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}

		start = time.Now()
		model, err := pdb.NewReader(file).Read()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(len(model.Atoms), time.Since(start))

		start = time.Now()
		mesh := ribbon.ModelMesh(model)
		fmt.Println(len(mesh.Triangles), time.Since(start))
	}
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "net/http/pprof"

	"github.com/fogleman/ribbon"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	var start time.Time
	for {
		start = time.Now()
		model, err := ribbon.LoadPDB(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(len(model.Atoms), time.Since(start))

		start = time.Now()
		mesh := model.Mesh()
		fmt.Println(len(mesh.Triangles), time.Since(start))
	}
}

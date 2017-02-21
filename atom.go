package ribbon

import "github.com/fogleman/fauxgl"

type Atom struct {
	Position   fauxgl.Vector
	Serial     int
	Name       string
	ResName    string
	ResSeq     int
	Occupancy  float64
	TempFactor float64
	Extra      string
}

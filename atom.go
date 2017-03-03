package ribbon

import "github.com/fogleman/fauxgl"

type Atom struct {
	Position fauxgl.Vector
	Name     string
	ResName  string
	ChainID  string
	ResSeq   int
	Element  string
	Het      bool
}

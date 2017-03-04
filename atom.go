package ribbon

import (
	"strings"

	"github.com/fogleman/fauxgl"
)

type Atom struct {
	Position fauxgl.Vector
	Serial   int
	Name     string
	ResName  string
	ChainID  string
	ResSeq   int
	Element  string
	Het      bool
}

func (a *Atom) GetElement() Element {
	return ElementsBySymbol[strings.Title(strings.ToLower(a.Element))]
}

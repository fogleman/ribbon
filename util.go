package ribbon

import (
	"strings"

	"github.com/fogleman/fauxgl"
	"github.com/fogleman/ribbon/pdb"
)

func atomPosition(a *pdb.Atom) fauxgl.Vector {
	return fauxgl.Vector{a.X, a.Y, a.Z}
}

func atomElement(a *pdb.Atom) Element {
	return ElementsBySymbol[strings.Title(strings.ToLower(a.Element))]
}

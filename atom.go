package ribbon

type Atom struct {
	Position   Vector
	Serial     int
	Name       string
	ResName    string
	ResSeq     int
	Occupancy  float64
	TempFactor float64
	Extra      string
}

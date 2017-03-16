package pdb

// COLUMNS        DATA  TYPE    FIELD        DEFINITION
// -------------------------------------------------------------------------------------
//  1 -  6        Record name   "ATOM  "
//  7 - 11        Integer       serial       Atom  serial number.
// 13 - 16        Atom          name         Atom name.
// 17             Character     altLoc       Alternate location indicator.
// 18 - 20        Residue name  resName      Residue name.
// 22             Character     chainID      Chain identifier.
// 23 - 26        Integer       resSeq       Residue sequence number.
// 27             AChar         iCode        Code for insertion of residues.
// 31 - 38        Real(8.3)     x            Orthogonal coordinates for X in Angstroms.
// 39 - 46        Real(8.3)     y            Orthogonal coordinates for Y in Angstroms.
// 47 - 54        Real(8.3)     z            Orthogonal coordinates for Z in Angstroms.
// 55 - 60        Real(6.2)     occupancy    Occupancy.
// 61 - 66        Real(6.2)     tempFactor   Temperature  factor.
// 77 - 78        LString(2)    element      Element symbol, right-justified.
// 79 - 80        LString(2)    charge       Charge  on the atom.

type Atom struct {
	Serial     int
	Name       string
	AltLoc     string
	ResName    string
	ChainID    string
	ResSeq     int
	ICode      string
	X, Y, Z    float64
	Occupancy  float64
	TempFactor float64
	Element    string
	Charge     string
}

func ParseAtom(line string) *Atom {
	atom := Atom{}
	atom.Serial = parseInt(line[6:11])
	atom.Name = parseString(line[12:16])
	atom.AltLoc = parseString(line[16:17])
	atom.ResName = parseString(line[17:20])
	atom.ChainID = parseString(line[21:22])
	atom.ResSeq = parseInt(line[22:26])
	atom.ICode = parseString(line[26:27])
	atom.Element = parseString(line[76:78])
	atom.X = parseFloat(line[30:38])
	atom.Y = parseFloat(line[38:46])
	atom.Z = parseFloat(line[46:54])
	atom.Occupancy = parseFloat(line[54:60])
	atom.TempFactor = parseFloat(line[60:66])
	atom.Charge = parseString(line[78:80])
	return &atom
}

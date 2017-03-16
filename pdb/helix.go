package pdb

// COLUMNS        DATA  TYPE     FIELD         DEFINITION
// -----------------------------------------------------------------------------------
//  1 -  6        Record name    "HELIX "
//  8 - 10        Integer        serNum        Serial number of the helix. This starts
//                                             at 1  and increases incrementally.
// 12 - 14        LString(3)     helixID       Helix  identifier. In addition to a serial
//                                             number, each helix is given an
//                                             alphanumeric character helix identifier.
// 16 - 18        Residue name   initResName   Name of the initial residue.
// 20             Character      initChainID   Chain identifier for the chain containing
//                                             this  helix.
// 22 - 25        Integer        initSeqNum    Sequence number of the initial residue.
// 26             AChar          initICode     Insertion code of the initial residue.
// 28 - 30        Residue  name  endResName    Name of the terminal residue of the helix.
// 32             Character      endChainID    Chain identifier for the chain containing
//                                             this  helix.
// 34 - 37        Integer        endSeqNum     Sequence number of the terminal residue.
// 38             AChar          endICode      Insertion code of the terminal residue.
// 39 - 40        Integer        helixClass    Helix class (see below).
// 41 - 70        String         comment       Comment about this helix.
// 72 - 76        Integer        length        Length of this helix.

type Helix struct {
	Serial      int
	HelixID     string
	InitResName string
	InitChainID string
	InitSeqNum  int
	InitICode   string
	EndResName  string
	EndChainID  string
	EndSeqNum   int
	EndICode    string
	HelixClass  int // TODO: enum?
	Length      int
}

func ParseHelix(line string) *Helix {
	helix := Helix{}
	helix.Serial = parseInt(line[7:10])
	helix.HelixID = parseString(line[11:14])
	helix.InitResName = parseString(line[15:18])
	helix.InitChainID = parseString(line[19:20])
	helix.InitSeqNum = parseInt(line[21:25])
	helix.InitICode = parseString(line[25:26])
	helix.EndResName = parseString(line[27:30])
	helix.EndChainID = parseString(line[31:32])
	helix.EndSeqNum = parseInt(line[33:37])
	helix.EndICode = parseString(line[37:38])
	helix.HelixClass = parseInt(line[38:40])
	helix.Length = parseInt(line[71:76])
	return &helix
}

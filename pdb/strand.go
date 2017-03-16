package pdb

// COLUMNS       DATA  TYPE     FIELD          DEFINITION
// -------------------------------------------------------------------------------------
//  1 -  6        Record name   "SHEET "
//  8 - 10        Integer       strand         Strand  number which starts at 1 for each
//                                             strand within a sheet and increases by one.
// 12 - 14        LString(3)    sheetID        Sheet  identifier.
// 15 - 16        Integer       numStrands     Number  of strands in sheet.
// 18 - 20        Residue name  initResName    Residue  name of initial residue.
// 22             Character     initChainID    Chain identifier of initial residue
//                                             in strand.
// 23 - 26        Integer       initSeqNum     Sequence number of initial residue
//                                             in strand.
// 27             AChar         initICode      Insertion code of initial residue
//                                             in  strand.
// 29 - 31        Residue name  endResName     Residue name of terminal residue.
// 33             Character     endChainID     Chain identifier of terminal residue.
// 34 - 37        Integer       endSeqNum      Sequence number of terminal residue.
// 38             AChar         endICode       Insertion code of terminal residue.
// 39 - 40        Integer       sense          Sense of strand with respect to previous
//                                             strand in the sheet. 0 if first strand,
//                                             1 if  parallel,and -1 if anti-parallel.
// 42 - 45        Atom          curAtom        Registration.  Atom name in current strand.
// 46 - 48        Residue name  curResName     Registration.  Residue name in current strand
// 50             Character     curChainId     Registration. Chain identifier in
//                                             current strand.
// 51 - 54        Integer       curResSeq      Registration.  Residue sequence number
//                                             in current strand.
// 55             AChar         curICode       Registration. Insertion code in
//                                             current strand.
// 57 - 60        Atom          prevAtom       Registration.  Atom name in previous strand.
// 61 - 63        Residue name  prevResName    Registration.  Residue name in
//                                             previous strand.
// 65             Character     prevChainId    Registration.  Chain identifier in
//                                             previous  strand.
// 66 - 69        Integer       prevResSeq     Registration. Residue sequence number
//                                             in previous strand.
// 70             AChar         prevICode      Registration.  Insertion code in
//                                             previous strand.

type Strand struct {
	Strand      int
	SheetID     string
	NumStrands  int
	InitResName string
	InitChainID string
	InitSeqNum  int
	InitICode   string
	EndResName  string
	EndChainID  string
	EndSeqNum   int
	EndICode    string
	Sense       int // TODO: enum?
	CurAtom     string
	CurResName  string
	CurChainId  string
	CurResSeq   int
	CurICode    string
	PrevAtom    string
	PrevResName string
	PrevChainId string
	PrevResSeq  int
	PrevICode   string
}

func ParseStrand(line string) *Strand {
	strand := Strand{}
	strand.Strand = parseInt(line[7:10])
	strand.SheetID = parseString(line[11:14])
	strand.NumStrands = parseInt(line[14:16])
	strand.InitResName = parseString(line[17:20])
	strand.InitChainID = parseString(line[21:22])
	strand.InitSeqNum = parseInt(line[22:26])
	strand.InitICode = parseString(line[26:27])
	strand.EndResName = parseString(line[28:31])
	strand.EndChainID = parseString(line[32:33])
	strand.EndSeqNum = parseInt(line[33:37])
	strand.EndICode = parseString(line[37:38])
	strand.Sense = parseInt(line[38:40])
	strand.CurAtom = parseString(line[41:45])
	strand.CurResName = parseString(line[45:48])
	strand.CurChainId = parseString(line[49:50])
	strand.CurResSeq = parseInt(line[50:54])
	strand.CurICode = parseString(line[54:55])
	strand.PrevAtom = parseString(line[56:60])
	strand.PrevResName = parseString(line[60:63])
	strand.PrevChainId = parseString(line[64:65])
	strand.PrevResSeq = parseInt(line[65:69])
	strand.PrevICode = parseString(line[69:70])
	return &strand
}

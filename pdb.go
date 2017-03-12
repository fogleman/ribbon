package ribbon

import (
	"bufio"
	"os"
	"strings"

	"github.com/fogleman/fauxgl"
)

func LoadPDB(path string) (*Model, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	var atoms []*Atom
	var hetAtoms []*Atom
	var connections []Connection
	var helixes []*Helix
	var strands []*Strand
	var bioMatrixes []fauxgl.Matrix
	var symMatrixes []fauxgl.Matrix
	var m [4][4]float64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "ENDMDL") {
			// TODO: handle multiple models
			break
		}
		if strings.HasPrefix(line, "REMARK 350   BIOMT") {
			row := parseInt(line[18:19]) - 1
			m[row][0] = parseFloat(strings.TrimSpace(line[23:33]))
			m[row][1] = parseFloat(strings.TrimSpace(line[33:43]))
			m[row][2] = parseFloat(strings.TrimSpace(line[43:53]))
			m[row][3] = parseFloat(strings.TrimSpace(line[53:68]))
			if row == 2 {
				bioMatrixes = append(bioMatrixes, fauxgl.Matrix{
					m[0][0], m[0][1], m[0][2], m[0][3],
					m[1][0], m[1][1], m[1][2], m[1][3],
					m[2][0], m[2][1], m[2][2], m[2][3],
					0, 0, 0, 1,
				})
			}
		}
		if strings.HasPrefix(line, "REMARK 290   SMTRY") {
			row := parseInt(line[18:19]) - 1
			m[row][0] = parseFloat(strings.TrimSpace(line[23:33]))
			m[row][1] = parseFloat(strings.TrimSpace(line[33:43]))
			m[row][2] = parseFloat(strings.TrimSpace(line[43:53]))
			m[row][3] = parseFloat(strings.TrimSpace(line[53:68]))
			if row == 2 {
				symMatrixes = append(symMatrixes, fauxgl.Matrix{
					m[0][0], m[0][1], m[0][2], m[0][3],
					m[1][0], m[1][1], m[1][2], m[1][3],
					m[2][0], m[2][1], m[2][2], m[2][3],
					0, 0, 0, 1,
				})
			}
		}
		if strings.HasPrefix(line, "ATOM  ") {
			atom := Atom{}
			x := parseFloat(strings.TrimSpace(line[30:38]))
			y := parseFloat(strings.TrimSpace(line[38:46]))
			z := parseFloat(strings.TrimSpace(line[46:54]))
			atom.Position = fauxgl.Vector{x, y, z}
			atom.Serial = parseInt(strings.TrimSpace(line[6:11]))
			atom.Name = strings.TrimSpace(line[12:16])
			atom.ResName = strings.TrimSpace(line[17:20])
			atom.ChainID = line[21:22]
			atom.ResSeq = parseInt(strings.TrimSpace(line[22:26]))
			atom.Element = strings.TrimSpace(line[76:78])
			atom.Occupancy = parseFloat(strings.TrimSpace(line[54:60]))
			atom.TempFactor = parseFloat(strings.TrimSpace(line[60:66]))
			atoms = append(atoms, &atom)
		}
		if strings.HasPrefix(line, "HETATM") {
			atom := Atom{}
			x := parseFloat(strings.TrimSpace(line[30:38]))
			y := parseFloat(strings.TrimSpace(line[38:46]))
			z := parseFloat(strings.TrimSpace(line[46:54]))
			atom.Position = fauxgl.Vector{x, y, z}
			atom.Serial = parseInt(strings.TrimSpace(line[6:11]))
			atom.Name = strings.TrimSpace(line[12:16])
			atom.ResName = strings.TrimSpace(line[17:20])
			atom.ChainID = line[21:22]
			atom.ResSeq = parseInt(strings.TrimSpace(line[22:26]))
			atom.Element = strings.TrimSpace(line[76:78])
			atom.Occupancy = parseFloat(strings.TrimSpace(line[54:60]))
			atom.TempFactor = parseFloat(strings.TrimSpace(line[60:66]))
			hetAtoms = append(hetAtoms, &atom)
		}
		if strings.HasPrefix(line, "CONECT") {
			a := parseInt(strings.TrimSpace(line[6:11]))
			b1 := parseInt(strings.TrimSpace(line[11:16]))
			b2 := parseInt(strings.TrimSpace(line[16:21]))
			b3 := parseInt(strings.TrimSpace(line[21:26]))
			b4 := parseInt(strings.TrimSpace(line[26:31]))
			for _, b := range []int{b1, b2, b3, b4} {
				if b != 0 {
					c := Connection{a, b}
					connections = append(connections, c)
				}
			}
		}
		if strings.HasPrefix(line, "HELIX ") {
			helix := Helix{}
			helix.ChainID = line[19:20]
			helix.InitSeqNum = parseInt(strings.TrimSpace(line[21:25]))
			helix.EndSeqNum = parseInt(strings.TrimSpace(line[33:37]))
			helixes = append(helixes, &helix)
		}
		if strings.HasPrefix(line, "SHEET ") {
			strand := Strand{}
			strand.ChainID = line[21:22]
			strand.InitSeqNum = parseInt(strings.TrimSpace(line[22:26]))
			strand.EndSeqNum = parseInt(strings.TrimSpace(line[33:37]))
			strands = append(strands, &strand)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	model := NewModel(atoms, helixes, strands)
	model.HetAtoms = hetAtoms
	model.Connections = connections
	model.BioMatrixes = bioMatrixes
	model.SymMatrixes = symMatrixes
	return model, nil
}

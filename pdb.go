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
	var helixes []*Helix
	var strands []*Strand
	var biologicalMatrixes []fauxgl.Matrix
	var symmetryMatrixes []fauxgl.Matrix
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
				biologicalMatrixes = append(biologicalMatrixes, fauxgl.Matrix{
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
				symmetryMatrixes = append(symmetryMatrixes, fauxgl.Matrix{
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
			atom.Chain = line[21:22]
			atom.ResSeq = parseInt(strings.TrimSpace(line[22:26]))
			atom.Occupancy = parseFloat(strings.TrimSpace(line[54:60]))
			atom.TempFactor = parseFloat(strings.TrimSpace(line[60:66]))
			atom.Element = strings.TrimSpace(line[76:78])
			atom.Extra = strings.TrimSpace(line[66:76])
			atoms = append(atoms, &atom)
		}
		if strings.HasPrefix(line, "HELIX ") {
			helix := Helix{}
			helix.Chain = line[19:20]
			helix.Start = parseInt(strings.TrimSpace(line[21:25]))
			helix.End = parseInt(strings.TrimSpace(line[33:37]))
			helixes = append(helixes, &helix)
		}
		if strings.HasPrefix(line, "SHEET ") {
			strand := Strand{}
			strand.Chain = line[21:22]
			strand.Start = parseInt(strings.TrimSpace(line[22:26]))
			strand.End = parseInt(strings.TrimSpace(line[33:37]))
			strands = append(strands, &strand)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	model := NewModel(atoms, helixes, strands)
	model.BiologicalMatrixes = biologicalMatrixes
	model.SymmetryMatrixes = symmetryMatrixes
	return model, nil
}

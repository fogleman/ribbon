package ribbon

import (
	"bufio"
	"os"
	"strings"
)

func LoadPDB(path string) (*Model, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	var atoms []*Atom
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "ATOM  ") {
			atom := Atom{}
			x := parseFloat(strings.TrimSpace(line[30:38]))
			y := parseFloat(strings.TrimSpace(line[38:46]))
			z := parseFloat(strings.TrimSpace(line[46:54]))
			atom.Position = Vector{x, y, z}
			atom.Serial = parseInt(strings.TrimSpace(line[6:11]))
			atom.Name = strings.TrimSpace(line[12:16])
			atom.ResName = strings.TrimSpace(line[17:20])
			atom.ResSeq = parseInt(strings.TrimSpace(line[22:26]))
			atom.Occupancy = parseFloat(strings.TrimSpace(line[54:60]))
			atom.TempFactor = parseFloat(strings.TrimSpace(line[60:66]))
			atom.Extra = strings.TrimSpace(line[66:76])
			atoms = append(atoms, &atom)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	model := NewModel(atoms)
	return model, nil
}

package pdb

import (
	"bufio"
	"io"
	"strings"
)

type Reader struct {
	io.Reader
}

func NewReader(r io.Reader) *Reader {
	return &Reader{r}
}

func (r *Reader) ReadAll() ([]*Model, error) {
	var models []*Model
	for {
		model, err := r.Read()
		if model != nil {
			models = append(models, model)
		}
		if err == io.EOF {
			return models, nil
		} else if err != nil {
			return models, err
		}
	}
}

func (r *Reader) Read() (*Model, error) {
	var ok bool
	var atoms []*Atom
	var hetAtoms []*Atom
	var connections []*Connection
	var helixes []*Helix
	var strands []*Strand
	var bioMatrixes []Matrix
	var symMatrixes []Matrix
	m := identity()
	scanner := bufio.NewScanner(r.Reader)
	for scanner.Scan() {
		line := scanner.Text()
		if ok && strings.HasPrefix(line, "ENDMDL") {
			break
		}
		if strings.HasPrefix(line, "ATOM  ") {
			atom := ParseAtom(line)
			atoms = append(atoms, atom)
			ok = true
		}
		if strings.HasPrefix(line, "HETATM") {
			atom := ParseAtom(line)
			hetAtoms = append(hetAtoms, atom)
			ok = true
		}
		if strings.HasPrefix(line, "CONECT") {
			cs := ParseConnections(line)
			connections = append(connections, cs...)
			ok = true
		}
		if strings.HasPrefix(line, "HELIX ") {
			helix := ParseHelix(line)
			helixes = append(helixes, helix)
			ok = true
		}
		if strings.HasPrefix(line, "SHEET ") {
			strand := ParseStrand(line)
			strands = append(strands, strand)
			ok = true
		}
		if strings.HasPrefix(line, "REMARK 350   BIOMT") {
			// TODO: per-chain matrices
			row := parseInt(line[18:19]) - 1
			m[row][0] = parseFloat(line[23:33])
			m[row][1] = parseFloat(line[33:43])
			m[row][2] = parseFloat(line[43:53])
			m[row][3] = parseFloat(line[53:68])
			if row == 2 {
				bioMatrixes = append(bioMatrixes, m)
				m = identity()
			}
		}
		if strings.HasPrefix(line, "REMARK 290   SMTRY") {
			row := parseInt(line[18:19]) - 1
			m[row][0] = parseFloat(line[23:33])
			m[row][1] = parseFloat(line[33:43])
			m[row][2] = parseFloat(line[43:53])
			m[row][3] = parseFloat(line[53:68])
			if row == 2 {
				symMatrixes = append(symMatrixes, m)
				m = identity()
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if !ok {
		return nil, io.EOF
	}
	residues := residuesForAtoms(atoms, helixes, strands)
	chains := chainsForResidues(residues)
	model := Model{}
	model.Atoms = atoms
	model.HetAtoms = hetAtoms
	model.Connections = connections
	model.Helixes = helixes
	model.Strands = strands
	model.BioMatrixes = bioMatrixes
	model.SymMatrixes = symMatrixes
	model.Residues = residues
	model.Chains = chains
	return &model, nil
}

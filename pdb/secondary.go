package pdb

type Secondary int

const (
	_ Secondary = iota
	SecondaryCoil
	SecondaryHelix
	SecondaryStrand
)

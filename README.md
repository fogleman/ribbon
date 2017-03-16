# Protein Ribbon Diagrams

Render ribbon diagrams of proteins in pure Go.

![4HHB](http://i.imgur.com/UFprBGt.png)

### Installation

[Go](https://golang.org/) should be installed and your `GOPATH` should be set (defaults to `$HOME/go` in Go 1.8+). `$GOPATH/bin` should be on your `$PATH` if you want to run the binaries easily.

    $ go get -u github.com/fogleman/ribbon/cmd/rcsb

### Example Usage

Provide a 4-digit RCSB Structure ID. The PDB file will automatically be downloaded and an image will be rendered. The triangle mesh will also be saved.

```bash
$ rcsb 4hhb  # generates 4hhb.png and 4hhb.stl
```

### TODO

Many features are not well exposed. It's not currently easy to configure things without changing code. This is what I hope to improve as I continue working on this project.

### Package `pdb`

The `pdb` package parses PDB files. The following entities are currently parsed:

```
ATOM   => *pdb.Atom
HETATM => *pdb.Atom
CONECT => *pdb.Connection
HELIX  => *pdb.Helix
SHEET  => *pdb.Strand
BIOMT  => pdb.Matrix
SMTRY  => pdb.Matrix
```

Additionally, some higher-level constructs are produced:

```
*pdb.Residue
*pdb.Chain
```

### Package `ribbon`

The `ribbon` package generates 3D meshes given a `pdb.Model`. It can produce the following types of meshes:

- Ribbon
- Ball & stick (for ligands)
- Space filling
- Backbone

### Package `fauxgl`

The [fauxgl](https://github.com/fogleman/fauxgl) library is used for rendering the 3D meshes in pure Go.

### Samples

![Sample](http://i.imgur.com/ImWjsrH.png)
![Sample](http://i.imgur.com/nQLRbfW.png)
![Sample](http://i.imgur.com/XNAgIoQ.png)
![Sample](http://i.imgur.com/YjQeClg.png)

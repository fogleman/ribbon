package ribbon

import (
	"math"

	"github.com/fogleman/fauxgl"
)

func ellipseProfile(n int, w, h float64) []fauxgl.Vector {
	result := make([]fauxgl.Vector, n)
	for i := range result {
		t := float64(i) / float64(n)
		a := t*2*math.Pi + math.Pi/4
		x := math.Cos(a) * w / 2
		y := math.Sin(a) * h / 2
		result[i] = fauxgl.Vector{x, y, 0}
	}
	return result
}

func rectangleProfile(n int, w, h float64) []fauxgl.Vector {
	result := make([]fauxgl.Vector, 0, n)
	hw := w / 2
	hh := h / 2
	segments := [][2]fauxgl.Vector{
		{fauxgl.Vector{hw, hh, 0}, fauxgl.Vector{-hw, hh, 0}},
		{fauxgl.Vector{-hw, hh, 0}, fauxgl.Vector{-hw, -hh, 0}},
		{fauxgl.Vector{-hw, -hh, 0}, fauxgl.Vector{hw, -hh, 0}},
		{fauxgl.Vector{hw, -hh, 0}, fauxgl.Vector{hw, hh, 0}},
	}
	m := n / 4
	for _, s := range segments {
		for i := 0; i < m; i++ {
			t := float64(i) / float64(m)
			p := s[0].Lerp(s[1], t)
			result = append(result, p)
		}
	}
	return result
}

func roundedRectangleProfile(n int, w, h float64) []fauxgl.Vector {
	result := make([]fauxgl.Vector, 0, n)
	r := h / 2
	hw := w/2 - r
	hh := h / 2
	segments := [][2]fauxgl.Vector{
		{fauxgl.Vector{hw, hh, 0}, fauxgl.Vector{-hw, hh, 0}},
		{fauxgl.Vector{-hw, 0, 0}, fauxgl.Vector{}},
		{fauxgl.Vector{-hw, -hh, 0}, fauxgl.Vector{hw, -hh, 0}},
		{fauxgl.Vector{hw, 0, 0}, fauxgl.Vector{}},
	}
	m := n / 4
	for si, s := range segments {
		for i := 0; i < m; i++ {
			t := float64(i) / float64(m)
			var p fauxgl.Vector
			switch si {
			case 0, 2:
				p = s[0].Lerp(s[1], t)
			case 1:
				a := math.Pi/2 + math.Pi*t
				x := math.Cos(a) * r
				y := math.Sin(a) * r
				p = s[0].Add(fauxgl.Vector{x, y, 0})
			case 3:
				a := 3*math.Pi/2 + math.Pi*t
				x := math.Cos(a) * r
				y := math.Sin(a) * r
				p = s[0].Add(fauxgl.Vector{x, y, 0})
			}
			result = append(result, p)
		}
	}
	return result
}

func translateProfile(p []fauxgl.Vector, dx, dy float64) []fauxgl.Vector {
	result := make([]fauxgl.Vector, len(p))
	for i := range result {
		result[i] = p[i].Add(fauxgl.Vector{dx, dy, 0})
	}
	return result
}

func geometryProfile(r1, r2 *Residue, n int) (p1, p2 []fauxgl.Vector) {
	switch r1.Type {
	case ResidueTypeHelix:
		p1 = roundedRectangleProfile(n, 3, 0.5)
		p1 = translateProfile(p1, 0, 1.5)
	case ResidueTypeStrand:
		if r2.Type == ResidueTypeStrand {
			p1 = rectangleProfile(n, 3, 1)
		} else {
			p1 = rectangleProfile(n, 4.5, 1)
		}
	default:
		p1 = ellipseProfile(n, 1, 1)
	}
	switch r2.Type {
	case ResidueTypeHelix:
		p2 = roundedRectangleProfile(n, 3, 0.5)
		p2 = translateProfile(p2, 0, 1.5)
	case ResidueTypeStrand:
		p2 = rectangleProfile(n, 3, 1)
	default:
		p2 = ellipseProfile(n, 1, 1)
	}
	return
}

func createSegmentMesh(pp1, pp2, pp3, pp4 *PeptidePlane) *fauxgl.Mesh {
	const splineSteps = 32
	const profileDetail = 32
	r1 := pp2.Residue1
	r2 := pp3.Residue1
	profile1, profile2 := geometryProfile(r1, r2, profileDetail)
	splines1 := make([][]fauxgl.Vector, len(profile1))
	splines2 := make([][]fauxgl.Vector, len(profile2))
	for i := range splines1 {
		p1 := profile1[i]
		p2 := profile2[i]
		splines1[i] = SplineForPlanes(pp1, pp2, pp3, pp4, splineSteps, p1.X, p1.Y)
		splines2[i] = SplineForPlanes(pp1, pp2, pp3, pp4, splineSteps, p2.X, p2.Y)
	}
	var triangles []*fauxgl.Triangle
	var lines []*fauxgl.Line
	for i := 0; i < splineSteps; i++ {
		t0 := float64(i) / splineSteps
		t1 := float64(i+1) / splineSteps
		if r2.Type == ResidueTypeStrand && r1.Type != ResidueTypeStrand {
			if t0 < 0.5 {
				t0 = 0
			} else {
				t0 = 1
			}
			if t1 < 0.5 {
				t1 = 0
			} else {
				t1 = 1
			}
		}
		for j := 0; j < profileDetail; j++ {
			p100 := splines1[j][i]
			p101 := splines1[j][i+1]
			p110 := splines1[(j+1)%profileDetail][i]
			p111 := splines1[(j+1)%profileDetail][i+1]
			p200 := splines2[j][i]
			p201 := splines2[j][i+1]
			p210 := splines2[(j+1)%profileDetail][i]
			p211 := splines2[(j+1)%profileDetail][i+1]
			p00 := p100.Lerp(p200, t0)
			p01 := p101.Lerp(p201, t1)
			p10 := p110.Lerp(p210, t0)
			p11 := p111.Lerp(p211, t1)
			triangles = triangulateQuad(triangles, p10, p11, p01, p00)
		}
	}
	return fauxgl.NewMesh(triangles, lines)
}

func triangulateQuad(triangles []*fauxgl.Triangle, p1, p2, p3, p4 fauxgl.Vector) []*fauxgl.Triangle {
	triangles = append(triangles, fauxgl.NewTriangleForPoints(p1, p2, p3))
	triangles = append(triangles, fauxgl.NewTriangleForPoints(p1, p3, p4))
	return triangles
}

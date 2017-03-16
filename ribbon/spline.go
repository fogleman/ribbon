package ribbon

import "github.com/fogleman/fauxgl"

func splineForPlanes(p1, p2, p3, p4 *PeptidePlane, n int, u, v float64) []fauxgl.Vector {
	g1 := p1.Position.Add(p1.Side.MulScalar(u)).Add(p1.Normal.MulScalar(v))
	g2 := p2.Position.Add(p2.Side.MulScalar(u)).Add(p2.Normal.MulScalar(v))
	g3 := p3.Position.Add(p3.Side.MulScalar(u)).Add(p3.Normal.MulScalar(v))
	g4 := p4.Position.Add(p4.Side.MulScalar(u)).Add(p4.Normal.MulScalar(v))
	return spline(g1, g2, g3, g4, n)
}

func spline(v1, v2, v3, v4 fauxgl.Vector, n int) []fauxgl.Vector {
	n1 := float64(n)
	n2 := float64(n * n)
	n3 := float64(n * n * n)
	s := fauxgl.Matrix{
		6 / n3, 0, 0, 0,
		6 / n3, 2 / n2, 0, 0,
		1 / n3, 1 / n2, 1 / n1, 0,
		0, 0, 0, 1,
	}
	b := fauxgl.Matrix{
		-1, 3, -3, 1,
		3, -6, 3, 0,
		-3, 0, 3, 0,
		1, 4, 1, 0,
	}.MulScalar(1.0 / 6.0)
	g := fauxgl.Matrix{
		v1.X, v1.Y, v1.Z, 1,
		v2.X, v2.Y, v2.Z, 1,
		v3.X, v3.Y, v3.Z, 1,
		v4.X, v4.Y, v4.Z, 1,
	}
	m := s.Mul(b).Mul(g)
	result := make([]fauxgl.Vector, 0, n+1)
	v := fauxgl.Vector{m.X30 / m.X33, m.X31 / m.X33, m.X32 / m.X33}
	v = v.RoundPlaces(10)
	result = append(result, v)
	for k := 0; k < n; k++ {
		m.X30 = m.X30 + m.X20
		m.X31 = m.X31 + m.X21
		m.X32 = m.X32 + m.X22
		m.X33 = m.X33 + m.X23
		m.X20 = m.X20 + m.X10
		m.X21 = m.X21 + m.X11
		m.X22 = m.X22 + m.X12
		m.X23 = m.X23 + m.X13
		m.X10 = m.X10 + m.X00
		m.X11 = m.X11 + m.X01
		m.X12 = m.X12 + m.X02
		m.X13 = m.X13 + m.X03
		v := fauxgl.Vector{m.X30 / m.X33, m.X31 / m.X33, m.X32 / m.X33}
		v = v.RoundPlaces(10)
		result = append(result, v)
	}
	return result
}

package ribbon

import (
	"math"

	"github.com/fogleman/fauxgl"
)

func OutlineSphere(eye, up, center fauxgl.Vector, radius float64) *fauxgl.Mesh {
	var lines []*fauxgl.Line
	hyp := center.Sub(eye).Length()
	opp := radius
	theta := math.Asin(opp / hyp)
	adj := opp / math.Tan(theta)
	d := math.Cos(theta) * adj
	r := math.Sin(theta) * adj
	w := center.Sub(eye).Normalize()
	u := w.Cross(up).Normalize()
	v := w.Cross(u).Normalize()
	c := eye.Add(w.MulScalar(d))
	var previous fauxgl.Vector
	for i := 0; i <= 360; i++ {
		a := fauxgl.Radians(float64(i))
		p := c
		p = p.Add(u.MulScalar(math.Cos(a) * r))
		p = p.Add(v.MulScalar(math.Sin(a) * r))
		if i > 0 {
			line := fauxgl.NewLineForPoints(previous, p)
			lines = append(lines, line)
		}
		previous = p
	}
	return fauxgl.NewLineMesh(lines)
}

func OutlineZCylinder(eye, up fauxgl.Vector, z0, z1, radius float64) *fauxgl.Mesh {
	center := fauxgl.Vector{0, 0, z0}
	hyp := center.Sub(eye).Length()
	opp := radius
	theta := math.Asin(opp / hyp)
	adj := opp / math.Tan(theta)
	d := math.Cos(theta) * adj
	w := center.Sub(eye).Normalize()
	u := w.Cross(up).Normalize()
	c0 := eye.Add(w.MulScalar(d))
	a0 := c0.Add(u.MulScalar(radius * 1.01))
	b0 := c0.Add(u.MulScalar(-radius * 1.01))

	center = fauxgl.Vector{0, 0, z1}
	hyp = center.Sub(eye).Length()
	opp = radius
	theta = math.Asin(opp / hyp)
	adj = opp / math.Tan(theta)
	d = math.Cos(theta) * adj
	w = center.Sub(eye).Normalize()
	u = w.Cross(up).Normalize()
	c1 := eye.Add(w.MulScalar(d))
	a1 := c1.Add(u.MulScalar(radius * 1.01))
	b1 := c1.Add(u.MulScalar(-radius * 1.01))

	var lines []*fauxgl.Line

	// for i := 0; i < 360; i++ {
	// 	a1 := fauxgl.Radians(float64(i))
	// 	a2 := fauxgl.Radians(float64(i + 1))
	// 	x1 := radius * math.Cos(a1)
	// 	y1 := radius * math.Sin(a1)
	// 	x2 := radius * math.Cos(a2)
	// 	y2 := radius * math.Sin(a2)
	// 	lines = append(lines, fauxgl.NewLineForPoints(
	// 		fauxgl.Vector{x1, y1, z0}, fauxgl.Vector{x2, y2, z0}))
	// 	lines = append(lines, fauxgl.NewLineForPoints(
	// 		fauxgl.Vector{x1, y1, z1}, fauxgl.Vector{x2, y2, z1}))
	// }

	for i := 0; i < 16; i++ {
		p1 := z0 + (z1-z0)*(float64(i)/16)
		p2 := z0 + (z1-z0)*(float64(i+1)/16)
		lines = append(lines, fauxgl.NewLineForPoints(
			fauxgl.Vector{a0.X, a0.Y, p1}, fauxgl.Vector{a1.X, a1.Y, p2}))
		lines = append(lines, fauxgl.NewLineForPoints(
			fauxgl.Vector{b0.X, b0.Y, p1}, fauxgl.Vector{b1.X, b1.Y, p2}))
	}

	return fauxgl.NewLineMesh(lines)
}

func OutlineCylinder(eye, up, v0, v1 fauxgl.Vector, radius float64) *fauxgl.Mesh {
	d := v1.Sub(v0)
	z := d.Length()
	matrix := fauxgl.RotateTo(fauxgl.Vector{0, 0, 1}, d.Normalize())
	matrix = matrix.Translate(v0)
	inverse := matrix.Inverse()
	mesh := OutlineZCylinder(
		inverse.MulPosition(eye), inverse.MulDirection(up), 0, z, radius)
	mesh.Transform(matrix)
	return mesh
}

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
	for q := 0; q < 3; q++ {
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
		r *= 1.01
	}
	return fauxgl.NewLineMesh(lines)
}

func OutlineCylinder(eye, v0, v1 fauxgl.Vector, radius float64) *fauxgl.Mesh {
	var lines []*fauxgl.Line
	for q := 0; q < 3; q++ {
		a := v1.Sub(v0).Normalize()
		b := v0.Sub(eye).Normalize()
		c := a.Cross(b).Normalize()
		a0 := v0.Add(c.MulScalar(radius))
		a1 := v1.Add(c.MulScalar(radius))
		b0 := v0.Add(c.MulScalar(-radius))
		b1 := v1.Add(c.MulScalar(-radius))
		const n = 36
		for i := 0; i < n; i++ {
			t0 := float64(i) / n
			t1 := float64(i+1) / n
			lines = append(lines, fauxgl.NewLineForPoints(a0.Lerp(a1, t0), a0.Lerp(a1, t1)))
			lines = append(lines, fauxgl.NewLineForPoints(b0.Lerp(b1, t0), b0.Lerp(b1, t1)))
		}
		radius *= 1.01
	}
	return fauxgl.NewLineMesh(lines)
}

func OutlineCylinderSphereIntersection(v0, v1 fauxgl.Vector, cr, sr float64) *fauxgl.Mesh {
	var lines []*fauxgl.Line
	z := math.Sqrt(sr*sr - cr*cr)
	d := v1.Sub(v0).Normalize()
	c := v0.Add(d.MulScalar(z))
	u := d.Perpendicular()
	v := u.Cross(d)
	var previous fauxgl.Vector
	for i := 0; i <= 360; i++ {
		a := fauxgl.Radians(float64(i))
		p := c
		p = p.Add(u.MulScalar(math.Cos(a) * cr))
		p = p.Add(v.MulScalar(math.Sin(a) * cr))
		if i > 0 {
			line := fauxgl.NewLineForPoints(previous, p)
			lines = append(lines, line)
		}
		previous = p
	}
	return fauxgl.NewLineMesh(lines)
}

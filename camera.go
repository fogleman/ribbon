package ribbon

import (
	"math"
	"math/rand"

	"github.com/fogleman/fauxgl"
)

type Camera struct {
	Eye    fauxgl.Vector
	Center fauxgl.Vector
	Up     fauxgl.Vector
	Fovy   float64
}

func makeCamera(points []fauxgl.Vector) Camera {
	const N = 500
	if len(points) > N {
		for i := range points {
			j := rand.Intn(i + 1)
			points[i], points[j] = points[j], points[i]
		}
		points = points[:N]
	}
	up := fauxgl.Vector{0, 0, 1}

	var center fauxgl.Vector
	// for _, point := range points {
	// 	center = center.Add(point)
	// }
	// center = center.DivScalar(float64(len(points)))

	var eye fauxgl.Vector
	best := 1e9
	for i := 0; i < 1000; i++ {
		v := fauxgl.RandomUnitVector().MulScalar(10)
		score := cameraScore(points, v)
		if score < best {
			best = score
			eye = v
		}
	}

	var fovy float64
	c := center.Sub(eye).Normalize()
	for _, point := range points {
		d := point.Sub(eye).Normalize()
		a := fauxgl.Degrees(math.Acos(d.Dot(c)) * 2 * 1.1)
		fovy = math.Max(fovy, a)
	}

	return Camera{eye, center, up, fovy}
}

func cameraScore(points []fauxgl.Vector, eye fauxgl.Vector) float64 {
	var result float64
	for _, p1 := range points {
		d1 := p1.Sub(eye).Normalize()
		for _, p2 := range points {
			d2 := p2.Sub(eye).Normalize()
			a := d1.Dot(d2)
			result += a * a
		}
	}
	return result
}

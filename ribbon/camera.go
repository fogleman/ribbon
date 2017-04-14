package ribbon

import (
	"math"
	"math/rand"

	"github.com/fogleman/fauxgl"
	"github.com/fogleman/ribbon/pdb"
)

type Camera struct {
	Eye    fauxgl.Vector
	Center fauxgl.Vector
	Up     fauxgl.Vector
	Fovy   float64
	Aspect float64
}

var DefaultCamera = Camera{
	fauxgl.Vector{0, 0, 5},
	fauxgl.Vector{0, 0, 0},
	fauxgl.Vector{0, 1, 0},
	30, 1,
}

func PositionCamera(model *pdb.Model) Camera {
	var points []fauxgl.Vector
	matrix := fauxgl.Identity()
	// for _, m := range model.SymMatrixes {
	// 	matrix := fauxgl.Matrix{
	// 		m[0][0], m[0][1], m[0][2], m[0][3],
	// 		m[1][0], m[1][1], m[1][2], m[1][3],
	// 		m[2][0], m[2][1], m[2][2], m[2][3],
	// 		m[3][0], m[3][1], m[3][2], m[3][3],
	// 	}
	for _, r := range model.Residues {
		if _, ok := r.AtomsByName["CA"]; !ok {
			continue
		}
		if _, ok := r.AtomsByName["O"]; !ok {
			continue
		}
		points = append(points, matrix.MulPosition(atomPosition(r.AtomsByName["CA"])))
		points = append(points, matrix.MulPosition(atomPosition(r.AtomsByName["O"])))
	}
	for _, a := range model.HetAtoms {
		points = append(points, matrix.MulPosition(atomPosition(a)))
	}
	// }
	if len(points) == 0 {
		return DefaultCamera
	}
	return makeCamera(points)
}

func makeCamera(points []fauxgl.Vector) Camera {
	const D = 1000
	up := fauxgl.Vector{0, 0, 1}

	min := points[0]
	max := points[0]
	for _, p := range points {
		min = min.Min(p)
		max = max.Max(p)
	}
	center := min.Add(max.Sub(min).MulScalar(0.5))

	_, r := bestBoundingSphere(points, 100)
	fovyEstimate := fauxgl.Degrees(math.Atan2(r, D) * 2.2)

	var eye fauxgl.Vector
	bestScore := 1e9
	size := int(math.Sqrt(float64(len(points))))
	for i := 0; i < 1000; i++ {
		v := fauxgl.RandomUnitVector().MulScalar(D)
		m := fauxgl.LookAt(v, center, up).Perspective(fovyEstimate, 1, 1, 100)
		score := cameraScore(points, m, size)
		if score < bestScore {
			bestScore = score
			eye = v
		}
	}

	bestAspect := 1.0
	bestFovy := fovyEstimate
	bestUp := up
	up = cameraUp(eye, center, up)
	forward := center.Sub(eye).Normalize()
	rotate := fauxgl.Rotate(forward, fauxgl.Radians(1))
	for i := 0; i < 180; i++ {
		m := fauxgl.LookAt(eye, center, up).Perspective(fovyEstimate, 1, 1, 100)
		w, h := cameraAspect(points, m)
		aspect := w / h
		if aspect > bestAspect {
			bestAspect = aspect
			bestFovy = fovyEstimate * h / 2
			bestUp = up
		}
		up = rotate.MulDirection(up)
	}
	aspect := bestAspect
	up = bestUp
	fovy := bestFovy * 1.1

	return Camera{eye, center, up, fovy, aspect}
}

func cameraUp(eye, center, up fauxgl.Vector) fauxgl.Vector {
	z := eye.Sub(center).Normalize()
	x := up.Cross(z).Normalize()
	y := z.Cross(x)
	return y
}

func cameraAspect(points []fauxgl.Vector, m fauxgl.Matrix) (float64, float64) {
	var w, h float64
	for _, p := range points {
		v := m.MulPositionW(p)
		v = v.DivScalar(v.W)
		w = math.Max(w, math.Abs(v.X)*2)
		h = math.Max(h, math.Abs(v.Y)*2)
	}
	return w, h
}

func cameraScore(points []fauxgl.Vector, m fauxgl.Matrix, size int) float64 {
	grid := make([]int, size*size)
	for _, p := range points {
		v := m.MulPositionW(p)
		v = v.DivScalar(v.W)
		x := int((v.X + 1) / 2 * float64(size))
		y := int((v.Y + 1) / 2 * float64(size))
		i := y*size + x
		if i >= 0 && i < len(grid) {
			grid[i]++
		}
	}
	var score float64
	for _, n := range grid {
		score += float64(n * n)
	}
	return score
}

func bestBoundingSphere(points []fauxgl.Vector, n int) (fauxgl.Vector, float64) {
	var minCenter fauxgl.Vector
	var minRadius float64
	for i := 0; i < n; i++ {
		c, r := boundingSphere(points)
		if i == 0 || r < minRadius {
			minCenter = c
			minRadius = r
		}
	}
	return minCenter, minRadius
}

func boundingSphere(points []fauxgl.Vector) (fauxgl.Vector, float64) {
	x := points[rand.Intn(len(points))]
	y, _ := furthestPoint(points, x)
	z, _ := furthestPoint(points, y)
	c := y.Lerp(z, 0.5)
	r := y.Distance(z) / 2
	_, r = furthestPoint(points, c)
	return c, r
}

func furthestPoint(points []fauxgl.Vector, point fauxgl.Vector) (fauxgl.Vector, float64) {
	var maxPoint fauxgl.Vector
	var maxDistance float64
	for i, p := range points {
		d := p.Distance(point)
		if i == 0 || d > maxDistance {
			maxPoint = p
			maxDistance = d
		}
	}
	return maxPoint, maxDistance
}

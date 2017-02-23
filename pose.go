package ribbon

import "github.com/fogleman/fauxgl"

type Pose struct {
	Position fauxgl.Vector
	Up       fauxgl.Vector
	Forward  fauxgl.Vector
	Right    fauxgl.Vector
}

func (pose *Pose) Point(p fauxgl.Vector) fauxgl.Vector {
	result := pose.Position
	result = result.Add(pose.Right.MulScalar(p.X))
	result = result.Add(pose.Up.MulScalar(p.Y))
	result = result.Add(pose.Forward.MulScalar(p.Z))
	return result
}

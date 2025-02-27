package util

import "math"

type Vec2 struct {
	X, Y float64
}

var (
	Vec2Right = Vec2{X: 1, Y: 0}
)

func NewVec2(x, y float64) Vec2 {
	return Vec2{
		X: x,
		Y: y,
	}
}

func V2Norm(x, y float64) (float64, float64) {
	l := math.Sqrt(x*x + y*y)
	if l == 0 {
		return 0, 0
	}

	return x / l, y / l
}

func V2Len(x, y float64) float64 {
	return math.Sqrt(x*x + y*y)
}

func V2Angle(x, y float64) float64 {
	return math.Atan2(y, x)
}

func (v Vec2) Rotate(angle float64) Vec2 {
	cos := math.Cos(angle)
	sin := -math.Sin(angle)
	return Vec2{
		X: v.X*cos - v.Y*sin,
		Y: v.X*sin + v.Y*cos,
	}
}

func (v Vec2) Mul(s float64) Vec2 {
	return Vec2{
		X: v.X * s,
		Y: v.Y * s,
	}
}

func (v Vec2) Add(v2 Vec2) Vec2 {
	return Vec2{
		X: v.X + v2.X,
		Y: v.Y + v2.Y,
	}
}

func (v Vec2) Values() (float64, float64) {
	return v.X, v.Y
}

func (v Vec2) Dot(v2 Vec2) float64 {
	return v.X*v2.X + v.Y*v2.Y
}

// LerpAngle interpolates between two angles a and b by a factor t.
func LerpAngle(a, b, t float64) float64 {
	diff := b - a
	if diff < -math.Pi {
		diff += 2 * math.Pi
	} else if diff > math.Pi {
		diff -= 2 * math.Pi
	}
	return a + diff*t
}

// ShortestAngleDirection returns the direction of the shortest path between two angles.
func ShortestAngleDirection(a, b float64) float64 {
	diff := b - a
	if diff > math.Pi {
		diff -= 2 * math.Pi
	} else if diff < -math.Pi {
		diff += 2 * math.Pi
	}
	return diff
}

func AngleNormalize(a float64) float64 {
	for a < 0 {
		a += 2 * math.Pi
	}
	for a >= 2*math.Pi {
		a -= 2 * math.Pi
	}
	return a
}

func DegToRad(deg float64) float64 {
	return deg * math.Pi / 180
}

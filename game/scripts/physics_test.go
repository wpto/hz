package scripts

import (
	"fmt"
	"testing"
)

func TestSolveCircleCollission(t *testing.T) {
	c := []struct {
		x1, y1, x2, y2, c1, c2 float64
	}{
		{0.0, 0.0, 0.01, 0.01, 10, 10},
	}

	for _, tc := range c {
		c1 := &PhysicsCircle{ID: 1, X: tc.x1, Y: tc.y1, Radius: tc.c1}
		c2 := &PhysicsCircle{ID: 2, X: tc.x2, Y: tc.y2, Radius: tc.c2}

		solveCircleCollision(c1, c2)
		fmt.Println(c1, c2)

	}
}

package scripts

import (
	"hz/game/scripts/physics"
	"hz/game/util"
)

// var rotationSpeed = util.DegToRad(360 * 3) // per second

func EnemyFollowTarget(p *physics.Physics, dt float64, e *Enemy, target EnemyTarget) {
	if util.V2Len(target.Walk.X-e.x, target.Walk.Y-e.y) < target.Radius {
		return
	}

	dx, dy := util.V2Norm(target.Walk.X-e.x, target.Walk.Y-e.y)

	// actualAngle := util.V2Angle(e.x-e.prevx, e.y-e.prevy)
	// targetAngle := util.V2Angle(targetX-e.x, targetY-e.y)

	// speed := util.ShortestAngleDirection(actualAngle, targetAngle)
	// if speed > rotationSpeed {
	// 	speed = rotationSpeed
	// }

	// stepAngle := actualAngle + speed*core.Delta

	// dx := math.Cos(stepAngle)
	// dy := math.Sin(stepAngle)

	e.SetLookDirection(target.Look)

	p.MoveAndCollide(e.id, dx*e.speed*dt, dy*e.speed*dt, 0)
	// e.SetPosition(e.x+dx*e.speed*dt, e.y+dy*e.speed*dt)
}

package scripts

import (
	"hz/game/core"
	"hz/game/util"
)

// var rotationSpeed = util.DegToRad(360 * 3) // per second

func EnemyFollowTarget(e *Enemy, target EnemyTarget) {
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

	e.SetPosition(e.x+dx*e.speed*core.Delta, e.y+dy*e.speed*core.Delta)
}

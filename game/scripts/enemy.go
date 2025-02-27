package scripts

import (
	"fmt"
	"hz/game/core"
	"hz/game/scripts/physics"
	"hz/game/util"
	"hz/resources/images"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Enemy struct {
	prevx, prevy float64
	x, y         float64
	speed        float64
	body         *core.Sprite
	legs         *core.AnimatedSprite

	lookDirection float64

	target EnemyTarget

	physics *physics.Physics
	id      int
}

func NewEnemy(p *physics.Physics, pos util.Vec2) *Enemy {
	enemy := &Enemy{
		x:       pos.X,
		y:       pos.Y,
		prevx:   pos.X,
		prevy:   pos.Y,
		speed:   120,
		body:    core.NewSprite(images.EnemyZombie),
		legs:    core.NewAnimatedSprite(images.Chel, "legs", 0.5),
		physics: p,
	}

	// id := p.AddCircle(0, 0, 10, enemy)
	enemy.id = p.AddShape(physics.CircleShape{X: 0, Y: 0, Radius: 10}, enemy)

	return enemy
}

func (e *Enemy) SetPosition(x, y float64) {
	e.prevx, e.prevy = e.x, e.y
	e.x, e.y = x, y
}

func (e *Enemy) SetLookDirection(target util.Vec2) {
	e.lookDirection = math.Atan2(target.Y-e.y, target.X-e.x)
}

func (e *Enemy) GetAnimationSpeed() float64 {
	return ((e.x-e.prevx)*(e.x-e.prevx) + (e.y-e.prevy)*(e.y-e.prevy)) / (e.speed) / e.speed / core.Delta / core.Delta
}

func (e *Enemy) PhysicsBackwardPush(shape physics.Shape) {
	if sh, ok := shape.(physics.CircleShape); ok {
		e.x, e.y = sh.X, sh.Y
	}
}

func (e *Enemy) PhysicsUpdate(dt float64) {
	if e.target.Active {
		EnemyFollowTarget(e.physics, dt, e, e.target)
	}
}

func (e *Enemy) Update() {
	fmt.Println("enemy position", e.x, e.y)
	e.body.SetPosition(e.x, e.y)
	e.body.SetRotation(e.lookDirection)
	e.body.Update()

	rotationForward := math.Atan2(e.y-e.prevy, e.x-e.prevx)
	e.legs.SetPosition(e.x, e.y)
	e.legs.SetRotation(rotationForward)
	e.legs.SetAnimationSpeed(e.GetAnimationSpeed())
	e.legs.Update()
}

func (e *Enemy) Draw(screen *ebiten.Image) {
	e.legs.Draw(screen)
	e.body.Draw(screen)
}

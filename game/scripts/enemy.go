package scripts

import (
	"hz/game/core"
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

	physics *Physics
	id      int
}

func NewEnemy(p *Physics, pos util.Vec2) *Enemy {
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

	id := p.AddCircle(0, 0, 10, enemy)
	enemy.id = id

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

func (e *Enemy) Update() {

	e.physics.UpdateCircle(e.id, e.x, e.y, 10)

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

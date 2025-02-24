package scripts

import (
	"hz/game/core"
	"hz/game/util"
	"hz/resources/images"

	"github.com/hajimehoshi/ebiten/v2"
)

type Bullet struct {
	x, y   float64
	dx, dy float64
	speed  float64

	sprite *core.Sprite
}

func NewBullet(x, y float64, dx, dy float64) *Bullet {
	return &Bullet{
		x:      x,
		y:      y,
		dx:     dx,
		dy:     dy,
		speed:  600,
		sprite: core.NewSprite(images.Bullet),
	}
}

func (b *Bullet) PhysicsUpdate(dt float64) {
	b.x += b.dx * b.speed * dt
	b.y += b.dy * b.speed * dt
}

func (b *Bullet) Update() {
	b.sprite.SetPosition(b.x, b.y)
	b.sprite.SetRotation(util.V2Angle(b.dx, b.dy))
	b.sprite.Update()
}

func (b *Bullet) Draw(screen *ebiten.Image) {
	b.sprite.Draw(screen)
}

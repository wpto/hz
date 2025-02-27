package scripts

import (
	"hz/game/scripts/physics"

	"github.com/hajimehoshi/ebiten/v2"
)

type BulletManager struct {
	Bullets []*Bullet

	Physics *physics.Physics
}

func NewBulletManager(ph *physics.Physics) *BulletManager {
	return &BulletManager{
		Physics: ph,
	}
}

func (bm *BulletManager) PhysicsUpdate(dt float64) {
	for _, b := range bm.Bullets {
		px, py := b.x, b.y
		b.PhysicsUpdate(dt)
		x, y := b.x, b.y

		_ = px
		_ = py
		_ = x
		_ = y
		// collide := bm.Physics.CheckLine(px, py, x, y)
		// if collide {
		// 	fmt.Println("Bullet collide")
		// }
	}
}

func (bm *BulletManager) Update() {
	for _, b := range bm.Bullets {
		b.Update()
	}
}

func (bm *BulletManager) Draw(screen *ebiten.Image) {
	for _, b := range bm.Bullets {
		b.Draw(screen)
	}
}

func (bm *BulletManager) AddBullet(b *Bullet) {
	bm.Bullets = append(bm.Bullets, b)
}

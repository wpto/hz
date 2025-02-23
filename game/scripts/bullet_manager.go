package scripts

import "github.com/hajimehoshi/ebiten/v2"

type BulletManager struct {
	Bullets []*Bullet
}

func NewBulletManager() *BulletManager {
	return &BulletManager{}
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
